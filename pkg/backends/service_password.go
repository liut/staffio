package backends

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dchest/passwordreset"
	"gopkg.in/mail.v2"

	"github.com/liut/staffio/pkg/common"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/settings"
)

var (
	ErrInvalidResetToken = errors.New("invalid reset token or not found")
	ErrMailNotReady      = errors.New("email system is not ready")
	ErrEmptyEmail        = errors.New("email is empty")

	secret []byte
)

func SetPasswordSecret(s string) {
	secret = []byte(s)
}

func (s *serviceImpl) getResetHash(uid string) ([]byte, error) {
	_, err := s.Get(uid)
	if err != nil {
		return nil, fmt.Errorf("no such user %s", uid)
	}
	uv, err := s.LoadVerify(uid)
	if err != nil {
		return nil, ErrInvalidResetToken
	}
	return uv.CodeHashBytes(), nil
}

func (s *serviceImpl) PasswordForgot(ctx context.Context, at common.AliasType, target, uid string) (err error) {
	var staff *models.Staff
	staff, err = s.Get(uid)
	if err != nil {
		return
	}
	if at != common.AtEmail {
		err = fmt.Errorf("invalid alias type %s", at.String())
		return
	}
	if at != common.AtEmail && target != staff.Email {
		err = fmt.Errorf("incorrect email %s", target)
		return
	}
	return s.passwordForgotPrepare(SiteFromContext(ctx), staff)
}

func (s *serviceImpl) passwordForgotPrepare(site string, staff *models.Staff) (err error) {
	if staff.Email == "" {
		return ErrEmptyEmail
	}
	uv := models.NewVerify(common.AtEmail, staff.Email, staff.UID)
	err = s.SaveVerify(uv)
	if err != nil {
		return
	}
	err = WriteUserLog(staff.UID, "password forgot", fmt.Sprintf("id %d, ch %d", uv.Id, uv.CodeHash))
	if err != nil {
		logger().Warnw("userLog fail", "uid", staff.UID, "err", err)
	}
	// Generate reset token that expires in 2 hours
	token := passwordreset.NewToken(staff.UID, 2*time.Hour, uv.CodeHashBytes(), secret)
	err = sendResetEmail(site, staff, token)
	return
}

func (s *serviceImpl) PasswordResetTokenVerify(token string) (uid string, err error) {
	uid, err = passwordreset.VerifyToken(token, s.getResetHash, secret)
	if err != nil {
		logger().Warnw("passwordreset.VerifyToken fail", "token", token, "err", err)
	}
	return
}

func (s *serviceImpl) PasswordResetWithToken(login, token, passwd string) (err error) {
	var uid string
	uid, err = s.PasswordResetTokenVerify(token)
	if err != nil {
		// verification failed, don't allow password reset
		return
	}
	if login != uid {
		return fmt.Errorf("invalid login %s", login)
	}
	// OK, reset password for uid (e.g. allow to change it)
	err = s.PasswordReset(uid, passwd)
	if err == nil {
		qs := func(db dbTxer) error {
			rs, de := db.Exec("DELETE FROM password_reset WHERE uid = $1", uid)
			if de == nil {
				ra, _ := rs.RowsAffected()
				logger().Infow("deleted reset", "affect", ra)
			}
			return de
		}
		err = withTxQuery(qs)
	}
	return
}

func (s *serviceImpl) SaveVerify(uv *models.Verify) error {
	qs := func(db dbTxer) error {
		euv, err := s.LoadVerify(uv.Uid)
		if err == nil {
			str := `DELETE FROM password_reset WHERE id = $1`
			_, err = db.Exec(str, euv.Id)
			if err != nil {
				logger().Warnw("DELETE password_reset fail", "uid", uv.Uid, "err", err)
				return err
			}
		}

		str := `INSERT INTO password_reset(type_id, target, uid, code_hash, life_seconds)
		 VALUES($1, $2, $3, $4, $5) RETURNING id`
		var id int
		err = db.Get(&id, str, uv.Type, uv.Target, uv.Uid, uv.CodeHash, uv.LifeSeconds)
		if err == nil {
			logger().Infow("new password_reset", "id", id, "uid", uv.Uid, "target", uv.Target)
			if id > 0 {
				uv.Id = id
			}

			return nil
		}
		logger().Warnw("INSERT password_reset fail", "uid", uv.Uid, "err", err)
		return err
	}
	return withTxQuery(qs)
}

func (s *serviceImpl) LoadVerify(uid string) (*models.Verify, error) {
	var uv models.Verify
	qs := func(db dber) error {
		return db.Get(&uv, `SELECT id, uid, type_id, target, code_hash, life_seconds, created, updated FROM password_reset
		 WHERE uid = $1 ORDER BY updated DESC LIMIT 1`, uid)
	}
	err := withDbQuery(qs)
	if err != nil {
		logger().Infow("query verify fail", "uid", uid, "err", err)
	}
	return &uv, err
}

var (
	replSite = strings.NewReplacer("www.", "i.")
)

func sendResetEmail(site string, staff *models.Staff, token string) error {
	var (
		smtpHost = settings.Current.MailHost
		smtpPort = settings.Current.MailPort
		smtpUser = settings.Current.MailSenderEmail
		smtpPass = settings.Current.MailSenderPassword
	)

	if !settings.Current.MailEnabled || smtpHost == "" {
		logger().Warnw("mail disabled or host is empty")
		return ErrMailNotReady
	}

	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", settings.Current.MailSenderName, smtpUser))
	m.SetHeader("To", staff.Email)
	m.SetHeader("Subject", "Password reset request")

	prefix := settings.Current.BaseURL
	switch site {
	case "i":
		prefix = replSite.Replace(prefix)
	case "www", "":
		prefix = prefix + "/password"
	}

	m.SetBody("text/html", fmt.Sprintf(tplPasswordReset, staff.Name(), prefix, token))

	logger().Infow("sending reset", "prefix", prefix, "email", staff.Email,
		"host", smtpHost, "port", smtpPort, "sender", smtpUser)

	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if settings.Current.MailTLSEnabled {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		logger().Infow("enable tls")
	}

	if err := d.DialAndSend(m); err != nil {
		logger().Warnw("send reset email failed", "host", smtpHost, "err", err)
		return err
	}
	logger().Infow("send reset email OK", "email", staff.Email)
	return nil
}

const (
	// tplPasswordReset = `Dear %s: <br/><br/>
	// To reset your password, pls <a href="%s/password/reset?rt=%s">click here</a>.`
	tplPasswordReset = `Dear %s: <br/><br/>
	To reset your password, pls <a href="%s/reset?token=%s">click here</a>.`
)

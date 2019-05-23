package backends

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dchest/passwordreset"
	"gopkg.in/mail.v2"

	"github.com/liut/staffio/pkg/common"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/settings"
)

var (
	ErrInvalidResetToken = errors.New("invalid reset token or not found")
	ErrEmptyMailhost     = errors.New("empty smtp host")
)

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

func (s *serviceImpl) PasswordForgot(at common.AliasType, target, uid string) (err error) {
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
	return s.passwordForgotPrepare(staff)
}

func (s *serviceImpl) passwordForgotPrepare(staff *models.Staff) (err error) {
	uv := models.NewVerify(common.AtEmail, staff.Email, staff.Uid)
	err = s.SaveVerify(uv)
	if err != nil {
		return
	}
	err = WriteUserLog(staff.Uid, "password forgot", fmt.Sprintf("id %d, ch %d", uv.Id, uv.CodeHash))
	if err != nil {
		log.Printf("userLog ERR %s", err)
	}
	// Generate reset token that expires in 2 hours
	secret := []byte(settings.PwdSecret)
	token := passwordreset.NewToken(staff.Uid, 2*time.Hour, uv.CodeHashBytes(), secret)
	err = sendResetEmail(staff, token)
	return
}

func (s *serviceImpl) PasswordResetTokenVerify(token string) (uid string, err error) {
	secret := []byte(settings.PwdSecret)
	uid, err = passwordreset.VerifyToken(token, s.getResetHash, secret)
	if err != nil {
		log.Printf("passwordreset.VerifyToken %q ERR %s", token, err)
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
				log.Printf("deleted %d", ra)
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
				log.Printf("DELETE password_reset %s ERR %s", uv.Uid, err)
				return err
			}
		}

		str := `INSERT INTO password_reset(type_id, target, uid, code_hash, life_seconds)
		 VALUES($1, $2, $3, $4, $5) RETURNING id`
		var id int
		err = db.Get(&id, str, uv.Type, uv.Target, uv.Uid, uv.CodeHash, uv.LifeSeconds)
		if err == nil {
			log.Printf("new password_reset id: %d of %s(%s)", id, uv.Uid, uv.Target)
			if id > 0 {
				uv.Id = id
			}

			return nil
		}
		log.Printf("INSERT password_reset %s ERR %s", uv.Uid, err)
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
		log.Printf("query verify with uid %q ERR %s", uid, err)
	}
	return &uv, err
}

func InitSMTP() {
	SetupSMTPhost(settings.SMTP.Host, settings.SMTP.Port)
	SetupSMTPAuth(settings.SMTP.SenderEmail, settings.SMTP.SenderPassword)
}

var (
	smtpHost string
	smtpPort = 587
	smtpUser string
	smtpPass string
)

func SetupSMTPhost(host string, port int) {
	smtpHost = host
	smtpPort = port
}

func SetupSMTPAuth(from, password string) {
	smtpUser = from
	smtpPass = password
}

func sendResetEmail(staff *models.Staff, token string) error {
	if smtpHost == "" {
		log.Print("smtp is disabled")
		return ErrEmptyMailhost
	}

	m := mail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", staff.Email)
	m.SetHeader("Subject", "Password reset request")
	m.SetBody("text/html", fmt.Sprintf(tplPasswordReset, staff.Name(), settings.BaseURL, token))

	logger().Infow("sending reset email", "email", staff.Email, "host", smtpHost)

	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

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

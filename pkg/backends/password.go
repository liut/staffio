package backends

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dchest/passwordreset"
	"github.com/wealthworks/csmtp"

	"lcgc/platform/staffio/pkg/models"
	"lcgc/platform/staffio/pkg/models/common"
	. "lcgc/platform/staffio/pkg/settings"
)

var (
	ErrInvalidResetToken = errors.New("invalid reset token or not found")
)

func getResetHash(uid string) ([]byte, error) {
	_, err := GetStaff(uid)
	if err != nil {
		return nil, fmt.Errorf("no such user %s", uid)
	}
	uv, err := LoadVerify(uid)
	if err != nil {
		return nil, ErrInvalidResetToken
	}
	return uv.CodeHashBytes(), nil
}

func PasswordForgot(at common.AliasType, target, uid string) (err error) {
	var staff *models.Staff
	staff, err = GetStaff(uid)
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
	return passwordForgotPrepare(staff)
}

func passwordForgotPrepare(staff *models.Staff) (err error) {
	uv := models.NewVerify(common.AtEmail, staff.Email, staff.Uid)
	err = SaveVerify(uv)
	if err != nil {
		return
	}
	err = WriteUserLog(staff.Uid, "password forgot", fmt.Sprintf("id %d, ch %d", uv.Id, uv.CodeHash))
	if err != nil {
		log.Printf("userLog ERR %s", err)
	}
	// Generate reset token that expires in 2 hours
	secret := []byte(Settings.PwdSecret)
	token := passwordreset.NewToken(staff.Uid, 2*time.Hour, uv.CodeHashBytes(), secret)
	err = sendResetEmail(staff, token)
	if err == nil {
		log.Printf("send reset email of %q OK", staff.Email)
	} else {
		log.Printf("send reset email ERR %s", err)
	}
	return
}

func PasswordResetTokenVerify(token string) (uid string, err error) {
	secret := []byte(Settings.PwdSecret)
	uid, err = passwordreset.VerifyToken(token, getResetHash, secret)
	if err != nil {
		log.Printf("passwordreset.VerifyToken %q ERR %s", token, err)
	}
	return
}

func PasswordResetWithToken(login, token, passwd string) (err error) {
	var uid string
	uid, err = PasswordResetTokenVerify(token)
	if err != nil {
		// verification failed, don't allow password reset
		return
	}
	if login != uid {
		return fmt.Errorf("invalid login %s", login)
	}
	// OK, reset password for uid (e.g. allow to change it)
	err = service.PasswordStore.PasswordReset(uid, passwd)
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

func SaveVerify(uv *models.Verify) error {
	qs := func(db dbTxer) error {
		log.Printf("save %v", uv)
		euv, err := LoadVerify(uv.Uid)
		if err == nil {
			str := `DELETE FROM password_reset WHERE id = $1`
			_, err = db.Exec(str, euv.Id)
			if err == nil {
				return nil
			}
			log.Printf("UPDATE password_reset %s ERR %s", uv.Uid, err)
			return err
		} else {
			log.Printf("DB ERR %s", err)
		}
		str := `INSERT INTO password_reset(type_id, target, uid, code_hash, life_seconds)
		 VALUES($1, $2, $3, $4, $5) RETURNING id`
		var id int
		err = db.Get(&id, str, uv.Type, uv.Target, uv.Uid, uv.CodeHash, uv.LifeSeconds)
		if err == nil {
			log.Printf("new password_reset id: %d of %s", id, uv.Uid)
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

func LoadVerify(uid string) (*models.Verify, error) {
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

func sendResetEmail(staff *models.Staff, token string) error {
	message := fmt.Sprintf(tplPasswordReset, staff.Name(), Settings.BaseURL, token)
	return csmtp.SendMail("Password reset request", message, staff.Email)
}

const (
	tplPasswordReset = `Dear %s:
	<br/><br/>
	To reset your password, pls <a href="%s/password/reset?rt=%s">click here</a>.`
)

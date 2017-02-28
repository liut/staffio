package backends

import (
	"fmt"
	"log"

	"lcgc/platform/staffio/backends/ldap"
	"lcgc/platform/staffio/models"
	. "lcgc/platform/staffio/settings"
)

var (
	backendReady bool
)

func Prepare() {
	if backendReady {
		return
	}

	ls := ldap.AddSource(Settings.LDAP.Host, Settings.LDAP.Base)

	ls.BindDN = Settings.LDAP.BindDN
	ls.Passwd = Settings.LDAP.Password
	ls.Debug = Settings.Debug

	backendReady = true
}

func CloseAll() {
	closeDb()
	ldap.CloseAll()
}

func LoadStaffs() []*models.Staff {
	limit := 20
	return ldap.ListPaged(limit)
}

func GetGroup(name string) *models.Group {
	return ldap.SearchGroup(name)
}

func GetStaff(uid string) (*models.Staff, error) {
	staff, err := ldap.GetStaff(uid)
	if err != nil {
		log.Printf("ldap get staff with %q ERR: %s", uid, err)
		return nil, err
	}
	return staff, nil
}

func StoreStaff(staff *models.Staff) error {
	isNew, err := ldap.StoreStaff(staff)
	if err == nil {
		if isNew {
			err = passwordForgotPrepare(staff)
			if err != nil {
				log.Printf("email of new user password send ERR %s", err)
			}
		}
	}
	return err
}

func DeleteStaff(uid string) error {
	return ldap.DeleteStaff(uid)
}

func Authenticate(uid, password string) bool {
	err := ldap.Authenticate(uid, password)
	if err != nil {
		log.Printf("Authen failed for %s, reason: %s", uid, err)
		return false
	}
	return true
}

func PasswordChange(uid, passwordOld, passwordNew string) error {
	err := ldap.PasswordChange(uid, passwordOld, passwordNew)
	if err != nil {
		if err == ldap.ErrLogin {
			return err
			// TODO:
		}
	}

	return err
}

func InGroup(group, uid string) bool {
	g := GetGroup(group)
	return g.Has(uid)
}

func ProfileModify(uid, password string, staff *models.Staff) error {
	if uid != staff.Uid {
		return fmt.Errorf("mismatch uid %s and %s", uid, staff.Uid)
	}
	return ldap.Modify(uid, password, staff)
}

func WriteUserLog(uid, subject, message string) error {
	qs := func(db dber) error {
		_, err := db.Exec("INSERT INTO user_log(uid, subject, body) VALUES($1, $2, $3)", uid, subject, message)
		return err
	}
	return withDbQuery(qs)
}

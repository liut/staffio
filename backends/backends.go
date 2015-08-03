package backends

import (
	"fmt"
	"log"
	"tuluu.com/liut/staffio/backends/ldap"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

var (
	backendReady bool
)

func Prepare() {
	if backendReady {
		return
	}
	addr := fmt.Sprintf("%s:%d", Settings.LDAP.Host, Settings.LDAP.Port)
	ls := ldap.AddSource(addr, Settings.LDAP.Base)
	ls.BindDN = Settings.LDAP.BindDN
	ls.Passwd = Settings.LDAP.Password
	ls.Debug = Settings.Debug

	backendReady = true
}

func CloseAll() {
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
		log.Printf("call GetStaff error: %s", err)
		return nil, err
	}
	return staff, nil
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

func ProfileModify(uid, password string, values map[string]string) error {
	// values["cn"] = fmt.Sprintf("%s%s", values["sn"], values["givenName"])
	return ldap.Modify(uid, password, values)
}

package backends

import (
	"errors"
	"fmt"
	"log"
	"tuluu.com/liut/staffio/backends/ldap"
	"tuluu.com/liut/staffio/models"
	. "tuluu.com/liut/staffio/settings"
)

var (
	ErrLogin = errors.New("Invalid Username/Password")
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

	backendReady = true
}

func CloseAll() {
	ldap.CloseAll()
}

func ListPaged(limit int) []*models.Staff {
	return ldap.ListPaged(limit)
}

func GetGroup(name string) *models.Group {
	return ldap.SearchGroup(name)
}

func GetStaff(uid string) (*models.Staff, error) {
	return ldap.GetStaff(uid)
}

func Authenticate(username, password string) (*models.Staff, error) {
	if ldap.Authenticate(username, password) {
		staff, err := ldap.GetStaff(username)
		if err != nil {
			log.Printf("call GetStaff error: %s", err)
		}
		return staff, nil
	}
	log.Printf("Login failed %s", username)
	return nil, ErrLogin
}

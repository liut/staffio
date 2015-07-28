package backends

import (
	"fmt"
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

	backendReady = true
}

func ListPaged(limit int) []*models.Staff {
	return ldap.ListPaged(limit)
}

func GetGroup(name string) *models.Group {
	return ldap.SearchGroup(name)
}

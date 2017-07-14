package backends

import (
	"log"

	"lcgc/platform/staffio/pkg/backends/ldap"
	"lcgc/platform/staffio/pkg/models"
	. "lcgc/platform/staffio/pkg/settings"
)

type Service struct {
	Authenticator models.Authenticator
	StaffStore    models.StaffStore
	PasswordStore models.PasswordStore
	GroupStore    models.GroupStore
	Close         func()
}

func NewService() *Service {

	cfg := &ldap.Config{
		Addr:   Settings.LDAP.Hosts,
		Base:   Settings.LDAP.Base,
		Bind:   Settings.LDAP.BindDN,
		Passwd: Settings.LDAP.Password,
	}
	store, err := ldap.NewStore(cfg)
	if err != nil {
		log.Fatalf("new service ERR %s", err)
	}
	// LDAP is a special store
	return &Service{
		Authenticator: store,
		StaffStore:    store,
		PasswordStore: store,
		GroupStore:    store,
		Close: func() {
			store.Close()
		},
	}

}

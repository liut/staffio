package ldap

import (
	"log"
	"strings"

	"github.com/liut/staffio/pkg/models"
)

var (
	_ models.Authenticator = (*LDAPStore)(nil)
	_ models.StaffStore    = (*LDAPStore)(nil)
	_ models.PasswordStore = (*LDAPStore)(nil)
	_ models.GroupStore    = (*LDAPStore)(nil)
)

type LDAPStore struct {
	sources  []*ldapSource
	pageSize int
}

func NewStore(cfg *Config) (*LDAPStore, error) {
	store := &LDAPStore{
		pageSize: 100,
	}
	for _, addr := range strings.Split(cfg.Addr, ",") {
		c := &Config{
			Addr:   addr,
			Base:   cfg.Base,
			Bind:   cfg.Bind,
			Passwd: cfg.Passwd,
		}
		ls, err := NewSource(c)
		if err != nil {
			return nil, err
		}
		store.sources = append(store.sources, ls)
	}

	return store, nil
}

func (s *LDAPStore) Authenticate(uid, passwd string) (err error) {
	for _, ls := range s.sources {
		dn := ls.UDN(uid)
		err = ls.Bind(dn, passwd, true)
		if err == nil {
			debug("authenticate(%s,****) ok", uid)
			return
		}
	}
	log.Printf("Authen failed for %s, reason: %s", uid, err)
	return
}

func (s *LDAPStore) Get(uid string) (staff *models.Staff, err error) {
	// log.Printf("sources %s", ldapSources)
	for _, ls := range s.sources {
		staff, err = ls.GetStaff(uid)
		if err == nil {
			return
		} else {
			log.Printf("GetStaff %s ERR %s", uid, err)
		}
	}
	err = ErrNotFound
	return
}

func (s *LDAPStore) All() (staffs []*models.Staff) {
	for _, ls := range s.sources {
		staffs = ls.ListPaged(s.pageSize)
		if len(staffs) > 0 {
			return
		}
	}
	return
}

func (s *LDAPStore) Save(staff *models.Staff) (isNew bool, err error) {
	for _, ls := range s.sources {
		isNew, err = ls.StoreStaff(staff)
		if err != nil {
			log.Printf("StoreStaff at %s ERR: %s", ls.Addr, err)
			return
		}
	}
	return
}

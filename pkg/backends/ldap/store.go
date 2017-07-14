package ldap

import (
	"log"
	"strings"

	"lcgc/platform/staffio/pkg/models"
)

var (
	_ models.Authenticator = (*storeImpl)(nil)
	_ models.StaffStore    = (*storeImpl)(nil)
	_ models.PasswordStore = (*storeImpl)(nil)
	_ models.GroupStore    = (*storeImpl)(nil)
)

type storeImpl struct {
	sources  []*ldapSource
	pageSize int
}

func NewStore(cfg *Config) (*storeImpl, error) {
	store := &storeImpl{
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

func (s *storeImpl) Authenticate(uid, passwd string) (err error) {
	for _, ls := range s.sources {
		dn := ls.UDN(uid)
		err = ls.Bind(dn, passwd, true)
		if err == nil {
			return nil
		}
	}
	return err
}

func (s *storeImpl) Get(uid string) (staff *models.Staff, err error) {
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

func (s *storeImpl) All() (staffs []*models.Staff) {
	for _, ls := range s.sources {
		staffs = ls.ListPaged(s.pageSize)
		if len(staffs) > 0 {
			return
		}
	}
	return
}

func (s *storeImpl) Save(staff *models.Staff) (isNew bool, err error) {
	for _, ls := range s.sources {
		isNew, err = ls.StoreStaff(staff)
		if err != nil {
			log.Printf("StoreStaff at %s ERR: %s", ls.Addr, err)
			return
		}
	}
	return
}

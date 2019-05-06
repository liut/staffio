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
		pageSize: PageSize,
	}
	for _, addr := range strings.Split(cfg.Addr, ",") {
		c := &Config{
			Addr:   addr,
			Base:   cfg.Base,
			Bind:   cfg.Bind,
			Passwd: cfg.Passwd,
		}
		ls, err := newSource(c)
		if err != nil {
			return nil, err
		}
		store.sources = append(store.sources, ls)
	}

	return store, nil
}

func (s *LDAPStore) Close() {
	for _, ls := range s.sources {
		ls.Close()
	}
}

func (s *LDAPStore) Authenticate(uid, passwd string) (err error) {
	for _, ls := range s.sources {
		err = ls.Authenticate(uid, passwd)
		if err == nil {
			debug("authenticate(%s,****) ok", uid)
			return
		}
	}
	log.Printf("Authen failed for %s, ERR: %s", uid, err)
	return
}

func (s *LDAPStore) Get(uid string) (staff *models.Staff, err error) {
	for _, ls := range s.sources {
		staff, err = ls.GetStaff(uid)
		if err == nil {
			return
		}
	}
	err = ErrNotFound
	return
}

func (s *LDAPStore) GetByDN(dn string) (staff *models.Staff, err error) {
	for _, ls := range s.sources {
		staff, err = ls.GetByDN(dn)
		if err == nil {
			return
		}
	}
	err = ErrNotFound
	return
}

func (s *LDAPStore) All() (staffs models.Staffs) {
	for _, ls := range s.sources {
		staffs = ls.ListPaged(s.pageSize)
	}
	return
}

func (s *LDAPStore) Save(staff *models.Staff) (isNew bool, err error) {
	for _, ls := range s.sources {
		isNew, err = ls.storeStaff(staff)
		if err != nil {
			log.Printf("storeStaff at %s ERR: %s", ls.Addr, err)
			return
		}
	}
	return
}

func (s *LDAPStore) Ready() error {
	for _, ls := range s.sources {
		err := ls.Ready("base", "groups", "people")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LDAPStore) PoolStats() *PoolStats {
	var pss PoolStats
	for _, ls := range s.sources {
		s := ls.cp.Stats()
		pss.Hits += s.Hits
		pss.Misses += s.Misses
		pss.Timeouts += s.Timeouts
		pss.TotalConns += s.TotalConns
		pss.IdleConns += s.IdleConns
	}
	return &pss
}
func splitDC(base string) string {
	pos1 := strings.Index(base, "=")
	pos2 := strings.Index(base, ",")
	// TODO:more condition
	return base[pos1+1 : pos2]
}

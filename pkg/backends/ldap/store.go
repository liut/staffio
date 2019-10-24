package ldap

import (
	"log"
	"strings"
)

type LDAPStore struct {
	sources  []*ldapSource
	pageSize int
}

func NewStore(cfg Config) (*LDAPStore, error) {
	if cfg.Base == "" {
		return nil, ErrEmptyBase
	}
	store := &LDAPStore{
		pageSize: cfg.PageSize,
	}
	for _, addr := range strings.Split(cfg.Addr, ",") {
		c := &Config{
			Addr:   addr,
			Base:   cfg.Base,
			Bind:   cfg.Bind,
			Passwd: cfg.Passwd,
			Domain: cfg.Domain,
		}
		ls, err := newSource(c)
		if err != nil {
			log.Printf("newSource(%s) ERR %s", addr, err)
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

// Authenticate verify uid and password from one of sources, return valid DN and error
func (s *LDAPStore) Authenticate(uid, passwd string) (staff *People, err error) {
	for _, ls := range s.sources {
		staff, err = ls.Authenticate(uid, passwd)
		if err == nil {
			debug("authenticate(%s,****) ok", uid)
			return
		}
	}
	log.Printf("Authen failed for %s, ERR: %s", uid, err)
	return
}

func (s *LDAPStore) Get(uid string) (staff *People, err error) {
	for _, ls := range s.sources {
		staff, err = ls.GetPeople(uid)
		if err == nil {
			return
		}
	}
	err = ErrNotFound
	return
}

func (s *LDAPStore) GetByDN(dn string) (staff *People, err error) {
	for _, ls := range s.sources {
		staff, err = ls.GetByDN(dn)
		if err == nil {
			return
		}
	}
	err = ErrNotFound
	return
}

func (s *LDAPStore) All(spec *Spec) (staffs Peoples) {
	if spec == nil {
		spec = new(Spec)
	}
	if spec.Limit == 0 {
		spec.Limit = s.pageSize
	}
	for _, ls := range s.sources {
		staffs = ls.List(spec)
	}
	return
}

func (s *LDAPStore) Save(staff *People) (isNew bool, err error) {
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

package ldap

import (
	"log"
	"strings"
)

// Store ..
type Store struct {
	sources  []*ldapSource
	pageSize int
}

// NewStore ...
func NewStore(cfg Config) (*Store, error) {
	if cfg.Base == "" {
		return nil, ErrEmptyBase
	}
	store := &Store{
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

// Close ...
func (s *Store) Close() {
	for _, ls := range s.sources {
		ls.Close()
	}
}

// Authenticate verify uid and password from one of sources, return valid DN and error
func (s *Store) Authenticate(uid, passwd string) (staff *People, err error) {
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

// Get return People with uid
func (s *Store) Get(uid string) (staff *People, err error) {
	for _, ls := range s.sources {
		staff, err = ls.GetPeople(uid)
		if err == nil {
			return
		}
	}
	err = ErrNotFound
	return
}

// GetByDN ...
func (s *Store) GetByDN(dn string) (staff *People, err error) {
	for _, ls := range s.sources {
		staff, err = ls.GetByDN(dn)
		if err == nil {
			return
		}
	}
	err = ErrNotFound
	return
}

// All ...
func (s *Store) All(spec *Spec) (staffs Peoples) {
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

// Save ...
func (s *Store) Save(staff *People) (isNew bool, err error) {
	for _, ls := range s.sources {
		isNew, err = ls.savePeople(staff)
		if err != nil {
			log.Printf("savePeople at %s ERR: %s", ls.Addr, err)
			return
		}
	}
	return
}

// Ready ...
func (s *Store) Ready() error {
	for _, ls := range s.sources {
		err := ls.Ready("base", "groups", "people")
		if err != nil {
			return err
		}
	}
	return nil
}

// PoolStats ...
func (s *Store) PoolStats() *PoolStats {
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

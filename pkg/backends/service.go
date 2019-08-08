package backends

import (
	"log"

	"github.com/liut/staffio/pkg/backends/ldap"
	"github.com/liut/staffio/pkg/common"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/cas"
	"github.com/liut/staffio/pkg/models/team"
	"github.com/liut/staffio/pkg/models/weekly"
)

var (
	ErrStoreNotFound = ldap.ErrNotFound
)

type PoolStats = ldap.PoolStats

type Servicer interface {
	models.Authenticator
	models.StaffStore
	models.PasswordStore
	models.GroupStore
	cas.TicketStore

	OSIN() OSINStore
	Ready() error
	CloseAll()

	SaveStaff(staff *models.Staff) error

	InGroup(gn, uid string) bool
	InGroupAny(uid string, names ...string) bool

	ProfileModify(uid, password string, staff *models.Staff) error

	PasswordForgot(at common.AliasType, target, uid string) error
	PasswordResetTokenVerify(token string) (uid string, err error)
	PasswordResetWithToken(login, token, passwd string) (err error)

	Team() team.Store
	Watch() team.WatchStore
	Weekly() weekly.Store

	PoolStats() *PoolStats
}

type serviceImpl struct {
	*ldap.LDAPStore
	osinStore   *DbStorage
	teamStore   *teamStore
	watchStore  *watchStore
	weeklyStore *weeklyStore
}

type LDAPConfig = ldap.Config

var ldapcfg *LDAPConfig

func SetLDAP(c LDAPConfig) {
	ldapcfg = &c
}

var _ Servicer = (*serviceImpl)(nil)

// NewService return new Servicer
func NewService() Servicer {
	cfg := ldap.NewConfig()
	if ldapcfg != nil {
		cfg.CopyFrom(*ldapcfg)
	}
	logger().Infow("new ldap config", "addr", cfg.Addr, "base", cfg.Base, "domain", cfg.Domain)

	store, err := ldap.NewStore(cfg)
	if err != nil {
		log.Fatalf("new service ERR %s", err)
	}
	// LDAP is a special store
	return &serviceImpl{
		LDAPStore:   store,
		osinStore:   NewStorage(),
		teamStore:   &teamStore{},
		watchStore:  &watchStore{store},
		weeklyStore: &weeklyStore{},
	}

}

func (s *serviceImpl) Ready() error {
	return s.LDAPStore.Ready()
}

func (s *serviceImpl) OSIN() OSINStore {
	return s.osinStore
}

func (s *serviceImpl) CloseAll() {
	s.LDAPStore.Close()
	s.osinStore.Close()
}

func (s *serviceImpl) Team() team.Store {
	return s.teamStore
}

func (s *serviceImpl) Watch() team.WatchStore {
	return s.watchStore
}

func (s *serviceImpl) Weekly() weekly.Store {
	return s.weeklyStore
}

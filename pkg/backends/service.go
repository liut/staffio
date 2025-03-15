package backends

import (
	"context"

	"github.com/liut/staffio-backend/ldap"
	schema "github.com/liut/staffio-backend/model"
	"github.com/liut/staffio/pkg/common"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/cas"
	"github.com/liut/staffio/pkg/models/team"
	"github.com/liut/staffio/pkg/models/weekly"
)

// vars
var (
	ErrStoreNotFound = ldap.ErrNotFound
)

// PoolStats ...
type PoolStats = ldap.PoolStats

// Group ...
type Group = schema.Group

// Spec ...
type Spec = schema.Spec

type PasswordStore interface {
	schema.PasswordStore
	PasswordForgot(ctx context.Context, at common.AliasType, target, uid string) error
	PasswordResetTokenVerify(token string) (uid string, err error)
	PasswordResetWithToken(login, token, passwd string) (err error)
}

// Servicer ...
type Servicer interface {
	schema.Authenticator
	schema.PeopleStore
	schema.GroupStore

	PasswordStore

	cas.TicketStore

	OSIN() OSINStore
	Ready() error
	CloseAll()

	SaveStaff(ctx context.Context, staff *models.Staff) error
	SaveTeamAndStaffs(ctx context.Context, team *team.Team, staffs models.Staffs) (err error)

	InGroup(gn, uid string) bool
	InGroupAny(uid string, names ...string) bool

	ProfileModify(uid, password string, staff *models.Staff) error

	Team() team.Store
	Watch() team.WatchStore
	Weekly() weekly.Store

	PoolStats() *PoolStats
}

type serviceImpl struct {
	*ldap.Store
	osinStore   *DbStorage
	teamStore   *teamStore
	watchStore  *watchStore
	weeklyStore *weeklyStore
}

// LDAPConfig ...
type LDAPConfig = ldap.Config

var ldapcfg *LDAPConfig

// SetLDAP ...
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
		logger().Fatalw("new service fail", "err", err)
	}
	// LDAP is a special store
	return &serviceImpl{
		Store:       store,
		osinStore:   NewStorage(),
		teamStore:   &teamStore{},
		watchStore:  &watchStore{store},
		weeklyStore: &weeklyStore{},
	}

}

func (s *serviceImpl) Ready() error {
	return s.Store.Ready()
}

func (s *serviceImpl) OSIN() OSINStore {
	return s.osinStore
}

func (s *serviceImpl) CloseAll() {
	s.Store.Close()
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

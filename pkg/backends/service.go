package backends

import (
	"log"

	"github.com/liut/staffio/pkg/backends/ldap"
	"github.com/liut/staffio/pkg/common"
	"github.com/liut/staffio/pkg/models"
	"github.com/liut/staffio/pkg/models/cas"
	"github.com/liut/staffio/pkg/models/weekly"
	"github.com/liut/staffio/pkg/settings"
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
	ProfileModify(uid, password string, staff *models.Staff) error
	PasswordForgot(at common.AliasType, target, uid string) error
	PasswordResetTokenVerify(token string) (uid string, err error)
	PasswordResetWithToken(login, token, passwd string) (err error)
	Team() weekly.TeamStore
	Weekly() weekly.WeeklyStore
	PoolStats() *PoolStats
}

type serviceImpl struct {
	*ldap.LDAPStore
	osinStore   *DbStorage
	teamStore   *teamStore
	weeklyStore *weeklyStore
}

var _ Servicer = (*serviceImpl)(nil)

// NewService return new Servicer
func NewService() Servicer {
	cfg := ldap.NewConfig()
	cfg.Addr = settings.LDAP.Hosts
	cfg.Base = settings.LDAP.Base
	cfg.Domain = settings.EmailDomain
	cfg.Bind = settings.LDAP.BindDN
	cfg.Passwd = settings.LDAP.Password
	store, err := ldap.NewStore(cfg)
	if err != nil {
		log.Fatalf("new service ERR %s", err)
	}
	// LDAP is a special store
	return &serviceImpl{
		LDAPStore:   store,
		osinStore:   NewStorage(),
		teamStore:   &teamStore{},
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

func (s *serviceImpl) Team() weekly.TeamStore {
	return s.teamStore
}

func (s *serviceImpl) Weekly() weekly.WeeklyStore {
	return s.weeklyStore
}

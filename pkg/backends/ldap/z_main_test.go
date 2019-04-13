package ldap

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liut/staffio/pkg/models"
)

const (
	defaultBase   = "dc=example,dc=org"
	defaultDomain = "example.org"
)

var (
	store *LDAPStore
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Ltime | log.Lshortfile)

	var err error
	Base = envOr("LDAP_BASE_DN", defaultBase)
	Domain = envOr("LDAP_DOMAIN", defaultDomain)

	cfg := &Config{
		Addr:   envOr("LDAP_ADDRS", "localhost"),
		Base:   Base,
		Bind:   envOr("LDAP_BIND_DN", "cn=admin,dc=example,dc=org"),
		Passwd: envOr("LDAP_PASSWD", "mypassword"),
	}
	store, err = NewStore(cfg)
	if err != nil {
		log.Fatalf("new store ERR %s", err)
	}
	err = store.Ready()
	if err != nil {
		log.Fatalf("store ready ERR %s", err)
	}
	defer store.Close()
	m.Run()
}

func TestSourceFailed(t *testing.T) {
	var err error
	_, err = NewStore(&Config{})
	assert.Error(t, err)
	assert.EqualError(t, err, ErrEmptyBase.Error())

	_, err = NewStore(&Config{Addr: ":bad", Base: defaultBase})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse")

	var s *ldapSource
	s, err = NewSource(&Config{
		Addr: "ldaps://localhost",
		Base: defaultBase,
	})
	assert.NoError(t, err)
	log.Printf("ldap source: %s", s)
	err = s.Ready()
	assert.Error(t, err)
	s.Close()
}

func TestStaffError(t *testing.T) {
	var err error
	_, err = store.Get("noexist")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrNotFound.Error())

	err = store.Delete("noexist")
	assert.Error(t, err)

	_, err = store.Save(&models.Staff{})
	assert.Error(t, err)
	_, err = store.Save(&models.Staff{Uid: "six"})
	assert.Error(t, err)

	err = store.Authenticate("baduid", "badPwd")
	assert.Error(t, err)
	assert.EqualError(t, err, ErrLogin.Error())
}

func TestStaff(t *testing.T) {
	var err error
	uid := "doe"
	cn := "doe"
	sn := "doe"
	password := "secret"
	staff := &models.Staff{
		Uid:        uid,
		CommonName: cn,
		Surname:    sn,
	}

	var isNew bool
	isNew, err = store.Save(staff)
	assert.NoError(t, err)
	assert.True(t, isNew)

	isNew, err = store.Save(staff)
	assert.NoError(t, err)
	assert.False(t, isNew)

	staff, err = store.Get(uid)
	assert.NoError(t, err)
	assert.Equal(t, cn, staff.CommonName)
	assert.Equal(t, sn, staff.Surname)

	data := store.All()
	assert.NotZero(t, len(data))

	err = store.PasswordReset(uid, password)
	assert.NoError(t, err)

	err = store.Authenticate(uid, password)
	assert.NoError(t, err)

	staff.GivenName = "fawn"
	err = store.ModifyBySelf(uid, password, staff)
	assert.NoError(t, err)

	err = store.PasswordChange(uid, password, "secretNew")
	assert.NoError(t, err)

	err = store.Delete(uid)
	assert.NoError(t, err)
}

func TestGroup(t *testing.T) {
	var err error
	_, err = store.GetGroup("noexist")
	assert.Error(t, err)
	_, err = store.AllGroup()
	assert.NoError(t, err)
}

func TestReady(t *testing.T) {
	var err error
	name := "teams"
	err = store.sources[0].Ready("")
	assert.NoError(t, err)
	err = store.sources[0].Ready(name)
	assert.NoError(t, err)

	err = store.sources[0].Delete(etParent.DN(name))
	assert.NoError(t, err)
}

package ldap

import (
	"log"
	"testing"
	"time"

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

func TestStoreFailed(t *testing.T) {
	var err error
	var _s *LDAPStore
	_, err = NewStore(&Config{})
	assert.Error(t, err)
	assert.EqualError(t, err, ErrEmptyBase.Error())

	_, err = NewStore(&Config{Addr: ":bad", Base: defaultBase})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse")

	_s, err = NewStore(&Config{
		Addr: "ldaps://localhost",
		Base: defaultBase,
	})
	assert.NoError(t, err)
	// log.Printf("ldap store: %s", _s)
	err = _s.Ready()
	assert.Error(t, err)
	_s.Close()
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
		// required fields
		Uid:        uid,
		CommonName: cn,
		Surname:    sn,

		// optional fields
		GivenName:      "fawn",
		AvatarPath:     "avatar.png",
		Description:    "It's me",
		Email:          "fawn@deer.cc",
		Nickname:       "tiny",
		Birthday:       "20120304",
		Gender:         models.Male,
		Mobile:         "13012341234",
		JoinDate:       time.Now().Format(DateLayout),
		EmployeeNumber: 001,
		EmployeeType:   "Engineer",
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

	staff.CommonName = "doe2"
	staff.GivenName = "fawn2"
	staff.Surname = "deer2"
	staff.AvatarPath = "avatar2.png"
	staff.Description = "It's me 2"
	staff.Email = "fawn2@deer.cc"
	staff.Nickname = "tiny2"
	staff.Birthday = "20120305"
	staff.Gender = models.Female
	staff.Mobile = "13012345678"
	staff.EmployeeNumber = 002
	staff.EmployeeType = "Chief Engineer"
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
	ls := store.sources[0]
	err = ls.Ready("")
	assert.NoError(t, err)
	err = ls.Ready(name)
	assert.NoError(t, err)

	err = ls.Delete(etParent.DN(name))
	assert.NoError(t, err)
}

func TestStoreStats(t *testing.T) {
	t.Logf("stats: %v", store.PoolStats())
}

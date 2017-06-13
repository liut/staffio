package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/go-ldap/ldap"

	"lcgc/platform/staffio/pkg/models"
)

var (
	_ models.Authenticator = (*storeImpl)(nil)
	_ models.StaffStore    = (*storeImpl)(nil)
	_ models.PasswordStore = (*storeImpl)(nil)
	_ models.GroupStore    = (*storeImpl)(nil)
)

// LDAP config
type Config struct {
	Addr, Base   string
	Bind, Passwd string
	Filter       string
	Attributes   []string
}

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

// Basic LDAP authentication service
type ldapSource struct {
	Addr       string     // LDAP address with host and port
	UseSSL     bool       // Use SSL
	Base       string     // Base DN
	BindDN     string     // default reader dn
	Passwd     string     // reader passwd
	Filter     string     // Query filter to validate entry
	Attributes []string   // Select fileds
	Enabled    bool       // if this source is disabled
	c          *ldap.Conn // conn
	bound      bool
	Debug      bool
}

var (
	ErrEmptyAddr      = errors.New("ldap addr is empty")
	ErrEmptyBase      = errors.New("ldap base is empty")
	ErrLogin          = errors.New("049: Invalid Username/Password")
	ErrNotFound       = errors.New("Not Found")
	userDnFmt         = "uid=%s,ou=people,%s"
	defaultFilter     = "(objectclass=inetOrgPerson)"
	defaultAttributes = []string{
		"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "description",
		"avatarPath", "dateOfBirth", "gender", "employeeNumber", "employeeType", "title"}
)

// Add a new source (LDAP directory) to the global pool
func NewSource(cfg *Config) (*ldapSource, error) {
	if cfg.Base == "" {
		return nil, ErrEmptyBase
	}

	log.Printf("new source %s", cfg.Addr)

	u, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("parse LDAP addr ERR: %s", err)
	}

	if u.Host == "" && u.Path != "" {
		u.Host = u.Path
		u.Path = ""
	}

	var useSSL bool
	if u.Scheme == "ldaps" {
		useSSL = true
	}

	pos := last(u.Host, ':')
	if pos < 0 {
		if useSSL {
			u.Host = u.Host + ":636"
		} else {
			u.Host = u.Host + ":389"
		}
	}

	filter := defaultFilter
	if cfg.Filter != "" {
		filter = cfg.Filter
	}

	ls := &ldapSource{
		Addr:       u.Host,
		UseSSL:     useSSL,
		Base:       cfg.Base,
		BindDN:     cfg.Bind,
		Passwd:     cfg.Passwd,
		Filter:     filter,
		Attributes: defaultAttributes,
		Enabled:    true,
	}

	return ls, nil
}

func (s *storeImpl) Close() {
	for _, ls := range s.sources {
		ls.Close()
	}
}

func (ls *ldapSource) String() string {
	return ls.Addr
}

func (ls *ldapSource) dial() (*ldap.Conn, error) {
	if ls.c != nil {
		return ls.c, nil
	}

	var err error
	if ls.UseSSL {
		ls.c, err = ldap.DialTLS("tcp", ls.Addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		ls.c, err = ldap.Dial("tcp", ls.Addr)
	}

	if err != nil {
		log.Printf("LDAP Connect error, %s:%v", ls.Addr, err)
		ls.Enabled = false
		return nil, err
	}

	// ls.c.Debug = ls.Debug

	return ls.c, nil
}

func (ls *ldapSource) Close() {
	ls.bound = false
	if ls.c != nil {
		ls.c.Close()
		ls.c = nil
	}
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

func (ls *ldapSource) UDN(uid string) string {
	return fmt.Sprintf(userDnFmt, uid, ls.Base)
}

func (ls *ldapSource) Bind(dn, passwd string, force bool) error {
	if !force && ls.bound {
		return nil
	}

	l, err := ls.dial()
	if err != nil {
		return err
	}

	err = l.Bind(dn, passwd)
	if err != nil {
		log.Printf("LDAP Bind failed for %s, reason: %s", dn, err.Error())
		if le, ok := err.(*ldap.Error); ok {
			if le.ResultCode == 49 {
				return ErrLogin
			}
		}
		return err
	}
	ls.bound = true
	return nil
}

func (ls *ldapSource) getEntry(udn string) (*ldap.Entry, error) {
	search := ldap.NewSearchRequest(
		udn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		ls.Filter,
		ls.Attributes,
		nil)
	sr, err := ls.c.Search(search)

	if err != nil {
		if le, ok := err.(*ldap.Error); ok && le.ResultCode == ldap.LDAPResultNoSuchObject {
			return nil, ErrNotFound
		}
		log.Printf("LDAP Search '%s' Error: %s", udn, err)
		return nil, err
	}

	if len(sr.Entries) > 0 {
		return sr.Entries[0], nil
	}
	return nil, ErrNotFound
}

// GetStaff : search an LDAP source if an entry (with uid) is valide and in the specific filter
func (ls *ldapSource) GetStaff(uid string) (*models.Staff, error) {
	err := ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		log.Printf("bind faild %s", err)
		return nil, err
	}

	entry, err := ls.getEntry(ls.UDN(uid))
	if err != nil {
		log.Printf("getEntry(%s) ERR %s", uid, err)
		return nil, err
	}

	return entryToUser(entry), nil
}

func (ls *ldapSource) ListPaged(limit int) (staffs []*models.Staff) {
	err := ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		// log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		return nil
	}

	if limit < 1 {
		limit = 1
	}
	search := ldap.NewSearchRequest(
		"ou=people,"+ls.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		ls.Filter,
		ls.Attributes,
		nil)

	sr, err := ls.c.SearchWithPaging(search, uint32(limit))
	if err != nil {
		log.Printf("ERROR: %s for search %v\n", err, search)
		return
	}

	if len(sr.Entries) > 0 {
		staffs = make([]*models.Staff, len(sr.Entries))
		for i, entry := range sr.Entries {
			staffs[i] = entryToUser(entry)
		}
	}

	return
}

func entryToUser(entry *ldap.Entry) (u *models.Staff) {
	// log.Printf("entry: %v", entry)
	u = new(models.Staff)
	u.Uid = entry.GetAttributeValue("uid")
	u.Surname = entry.GetAttributeValue("sn")
	u.GivenName = entry.GetAttributeValue("givenName")
	u.CommonName = entry.GetAttributeValue("cn")
	u.Email = entry.GetAttributeValue("mail")
	u.Nickname = entry.GetAttributeValue("displayName")
	u.Mobile = entry.GetAttributeValue("mobile")
	u.EmployeeNumber = entry.GetAttributeValue("employeeNumber")
	u.EmployeeType = entry.GetAttributeValue("employeeType")
	u.Birthday = entry.GetAttributeValue("dateOfBirth")
	(&u.Gender).UnmarshalJSON([]byte(entry.GetAttributeValue("gender")))
	u.AvatarPath = entry.GetAttributeValue("avatarPath")
	u.Description = entry.GetAttributeValue("description")
	return
}

// Index of rightmost occurrence of b in s.
func last(s string, b byte) int {
	i := len(s)
	for i--; i >= 0; i-- {
		if s[i] == b {
			break
		}
	}
	return i
}

package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/go-ldap/ldap"

	"lcgc/platform/staffio/models"
)

// Basic LDAP authentication service
type LdapSource struct {
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

//Global LDAP directory pool
var (
	userDnFmt         = "uid=%s,ou=people,%s"
	defaultAttributes = []string{"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "employeeNumber", "employeeType", "description", "title"}
	defaultFilter     = "(objectclass=inetOrgPerson)"
	AuthenSource      []*LdapSource
	ErrLogin          = errors.New("049: Invalid Username/Password")
	ErrNotFound       = errors.New("Not Found")
)

// Add a new source (LDAP directory) to the global pool
func AddSource(addr, base string) *LdapSource {
	if base == "" {
		log.Fatal("ldap base is empty")
	}

	u, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("parse LDAP Host ERR: %s", err)
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

	ls := &LdapSource{
		Addr:       u.Host,
		UseSSL:     useSSL,
		Base:       base,
		Filter:     defaultFilter,
		Attributes: defaultAttributes,
		Enabled:    true,
	}

	AuthenSource = append(AuthenSource, ls)
	return ls
}

func CloseAll() {
	for _, ls := range AuthenSource {
		ls.Close()
	}
}

func (ls *LdapSource) dial() (*ldap.Conn, error) {
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

func (ls *LdapSource) Close() {
	ls.bound = false
	if ls.c != nil {
		ls.c.Close()
		ls.c = nil
	}
}

func Authenticate(uid, passwd string) (err error) {
	for _, ls := range AuthenSource {
		dn := ls.UDN(uid)
		err = ls.Bind(dn, passwd, true)
		if err == nil {
			return nil
		}
	}
	return err
}

func GetStaff(uid string) (staff *models.Staff, err error) {
	for _, ls := range AuthenSource {
		staff, err = ls.GetStaff(uid)
		if err == nil {
			return
		}
	}
	return
}

func ListPaged(limit int) (staffs []*models.Staff) {
	for _, ls := range AuthenSource {
		staffs = ls.ListPaged(limit)
		if len(staffs) > 0 {
			return
		}
	}
	return
}

func (ls *LdapSource) UDN(uid string) string {
	return fmt.Sprintf(userDnFmt, uid, ls.Base)
}

func (ls *LdapSource) Bind(dn, passwd string, force bool) error {
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

func (ls *LdapSource) getEntry(udn string) (*ldap.Entry, error) {
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
		log.Printf("LDAP Search '%s' Error: ", udn, err)
		return nil, err
	}

	if len(sr.Entries) > 0 {
		return sr.Entries[0], nil
	}
	return nil, ErrNotFound
}

// GetStaff : search an LDAP source if an entry (with uid) is valide and in the specific filter
func (ls *LdapSource) GetStaff(uid string) (*models.Staff, error) {
	err := ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		return nil, err
	}

	entry, err := ls.getEntry(ls.UDN(uid))
	if err != nil {
		return nil, err
	}

	return entryToUser(entry), nil
}

func (ls *LdapSource) ListPaged(limit int) (staffs []*models.Staff) {
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

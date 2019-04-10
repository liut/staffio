package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/go-ldap/ldap"
	. "github.com/wealthworks/go-debug"

	"github.com/liut/staffio/pkg/models"
)

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
}

var (
	ErrEmptyAddr = errors.New("ldap addr is empty")
	ErrEmptyBase = errors.New("ldap base is empty")
	ErrLogin     = errors.New("049: Invalid Username/Password")
	ErrNotFound  = errors.New("Not Found")
	userDnFmt    = "uid=%s,ou=people,%s"

	debug = Debug("staffio:ldap")
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

	filter := etPeople.Filter
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
		Attributes: etPeople.Attributes,
		Enabled:    true,
	}

	return ls, nil
}

func (s *LDAPStore) Close() {
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

	debug("connect to %s ok", ls.Addr)
	return ls.c, nil
}

func (ls *ldapSource) Close() {
	ls.bound = false
	if ls.c != nil {
		ls.c.Close()
		ls.c = nil
	}
}

func (ls *ldapSource) UDN(uid string) string {
	return etPeople.DN(uid)
}

func (ls *ldapSource) Ready() (err error) {
	err = ls.Bind(ls.BindDN, ls.Passwd, false)
	if err == nil {
		if err = ls.readyBase(); err != nil {
			return
		}
		err = ls.readyParent("groups")
		if err == nil {
			err = ls.readyParent("people")
		}
	}
	return
}

func (ls *ldapSource) readyBase() (err error) {
	dn := Base
	_, err = ls.Entry(dn, etBase.Filter, etBase.Attributes...)
	if err == ErrNotFound {
		ar := ldap.NewAddRequest(dn, nil)
		ar.Attribute("objectClass", []string{etBase.OC, "organization", "top"})
		ar.Attribute("o", []string{Domain})
		ar.Attribute(etBase.PK, []string{splitDC(Base)})
		debug("add %v", ar)
		err = ls.c.Add(ar)
		if err != nil {
			debug("add %q, ERR: %s", dn, err)
		} else {
			debug("add %q OK", dn)
		}
	}
	return
}

func (ls *ldapSource) readyParent(name string) (err error) {
	dn := etParent.DN(name)
	_, err = ls.Entry(dn, etParent.Filter, etParent.Attributes...)
	if err == ErrNotFound {
		debug("ready parent %s, ERR %s", name, err)
		ar := ldap.NewAddRequest(dn, nil)
		ar.Attribute("objectClass", []string{etParent.OC, "top"})
		ar.Attribute(etParent.PK, []string{name})
		err = ls.c.Add(ar)
		if err != nil {
			debug("add %q, ERR: %s", dn, err)
		}
	}
	return
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

	debug("bind(%s, ***) ok", dn)
	ls.bound = true
	return nil
}

// deprecated with Entry(dn, filter string, attrs ...string)
func (ls *ldapSource) getEntry(udn string) (*ldap.Entry, error) {
	return ls.Entry(udn, ls.Filter, ls.Attributes...)
}

func (ls *ldapSource) Group(cn string) (*ldap.Entry, error) {
	return ls.Entry(ls.GDN(cn), etGroup.Filter, etGroup.Attributes...)
}

func (ls *ldapSource) People(uid string) (*ldap.Entry, error) {
	return ls.Entry(ls.UDN(uid), ls.Filter, ls.Attributes...)
}

// Entry return a special entry with dn and filter
func (ls *ldapSource) Entry(dn, filter string, attrs ...string) (*ldap.Entry, error) {
	if !ls.bound {
		err := ls.Bind(ls.BindDN, ls.Passwd, false)
		if err != nil {
			return nil, err
		}
	}
	search := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attrs,
		nil)
	sr, err := ls.c.Search(search)
	if err == nil {
		if len(sr.Entries) > 0 {
			return sr.Entries[0], nil
		}
		return nil, ErrNotFound
	}

	debug("ldap search %q, ERR: %s", dn, err)
	if le, ok := err.(*ldap.Error); ok && le.ResultCode == ldap.LDAPResultNoSuchObject {
		return nil, ErrNotFound
	}
	log.Printf("LDAP Search '%s' Error: %s", dn, err)
	return nil, err

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

func (ls *ldapSource) ListPaged(limit int) (staffs models.Staffs) {
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
		staffs = make(models.Staffs, len(sr.Entries))
		for i, entry := range sr.Entries {
			staffs[i] = entryToUser(entry)
		}
	}

	return
}

func entryToUser(entry *ldap.Entry) (u *models.Staff) {
	// log.Printf("entry: %v", entry)
	u = &models.Staff{
		Uid:          entry.GetAttributeValue("uid"),
		Surname:      entry.GetAttributeValue("sn"),
		GivenName:    entry.GetAttributeValue("givenName"),
		CommonName:   entry.GetAttributeValue("cn"),
		Email:        entry.GetAttributeValue("mail"),
		Nickname:     entry.GetAttributeValue("displayName"),
		Mobile:       entry.GetAttributeValue("mobile"),
		EmployeeType: entry.GetAttributeValue("employeeType"),
		Birthday:     entry.GetAttributeValue("dateOfBirth"),
		AvatarPath:   entry.GetAttributeValue("avatarPath"),
		Description:  entry.GetAttributeValue("description"),
		JoinDate:     entry.GetAttributeValue("dateOfJoin"),
		IDCN:         entry.GetAttributeValue("idcnNumber"),
	}
	(&u.Gender).UnmarshalText(entry.GetRawAttributeValue("gender"))
	var err error
	if str := entry.GetAttributeValue("employeeNumber"); str != "" {
		u.EmployeeNumber, err = strconv.Atoi(str)
		if err != nil {
			log.Printf("invalid employee number %q, ERR %s", str, err)
		}
	}
	if str := entry.GetAttributeValue("createdTime"); str != "" {
		u.Created, err = time.Parse(timeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	} else if str := entry.GetAttributeValue("createTimestamp"); str != "" {
		u.Created, err = time.Parse(timeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	}
	if str := entry.GetAttributeValue("modifiedTime"); str != "" {
		u.Updated, err = time.Parse(timeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	} else if str := entry.GetAttributeValue("modifyTimestamp"); str != "" {
		u.Updated, err = time.Parse(timeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	}
	if blob := entry.GetRawAttributeValue("jpegPhoto"); len(blob) > 0 {
		u.JpegPhoto = blob
	}
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

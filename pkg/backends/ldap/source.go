package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-ldap/ldap"
	. "github.com/wealthworks/go-debug"

	"github.com/liut/staffio/pkg/backends/ldap/pool"
	"github.com/liut/staffio/pkg/models"
)

type PoolStats = pool.Stats

// Basic LDAP authentication service
type ldapSource struct {
	Addr   string      // LDAP address with host and port
	Base   string      // Base DN
	BindDN string      // default reader dn
	Passwd string      // reader passwd
	cp     pool.Pooler // conn
}

var (
	ErrEmptyAddr = errors.New("ldap addr is empty")
	ErrEmptyBase = errors.New("ldap base is empty")
	ErrLogin     = errors.New("049: Invalid Username/Password")
	ErrNotFound  = errors.New("Not Found")
	userDnFmt    = "uid=%s,ou=people,%s"

	once sync.Once

	debug = Debug("staffio:ldap")
)

// newSource Add a new source (LDAP directory) to the global pool
func newSource(cfg *Config) (*ldapSource, error) {
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

	opt := &pool.Options{
		Factory: func() (ldap.Client, error) {
			if useSSL {
				return ldap.DialTLS("tcp", u.Host, &tls.Config{InsecureSkipVerify: true})
			}
			return ldap.Dial("tcp", u.Host)
		},
		PoolSize:           PoolSize,
		PoolTimeout:        30 * time.Second,
		MaxConnAge:         25 * time.Minute,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 2 * time.Minute,
	}

	ls := &ldapSource{
		Addr:   u.Host,
		Base:   cfg.Base,
		BindDN: cfg.Bind,
		Passwd: cfg.Passwd,
		cp:     pool.NewPool(opt),
	}

	return ls, nil
}

func (ls *ldapSource) Close() {
	if ls.cp != nil {
		ls.cp.Close()
	}
}

func (ls *ldapSource) UDN(uid string) string {
	if isADsource {
		etADuser.DN(uid)
	}
	return etPeople.DN(uid)
}

func (ls *ldapSource) Ready(names ...string) (err error) {
	err = ls.opWithMan(func(c ldap.Client) (err error) {
		for _, name := range names {
			if name == "" {
				continue
			}
			if name == "base" {
				err = ldapEntryReady(c, etBase, splitDC(ls.Base))
			} else {
				err = ldapEntryReady(c, etParent, name)
			}
		}
		return
	})
	return
}

func ldapEntryReady(c ldap.Client, et *entryType, name string) (err error) {
	dn := et.DN(name)
	var entry *ldap.Entry
	entry, err = ldapEntryGet(c, dn, et.Filter, et.Attributes...)
	debug("check ready for %s done, ERR %v", name, err)
	if err == nil && et == etBase {
		once.Do(func() {
			if entry.GetAttributeValue("instanceType") != "" {
				debug("The source is Active Directory!")
				isADsource = true
			}
		})
		return
	}
	if err == ErrNotFound {
		ar := ldap.NewAddRequest(dn, nil)
		// ar.Attribute("objectClass", []string{et.OC, "top"})
		// ar.Attribute(et.PK, []string{name})
		et.prepareTo(name, ar)
		debug("add %v", ar)
		err = c.Add(ar)
		if err != nil {
			debug("add %q, ERR: %s", dn, err)
		} else {
			debug("add %q OK", dn)
		}
		return
	}
	return
}

func ldapEntryGet(c ldap.Client, dn, filter string, attrs ...string) (*ldap.Entry, error) {
	search := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attrs,
		nil)
	sr, err := c.Search(search)
	if err == nil {
		if len(sr.Entries) > 0 {
			debug("found dn %q entries %d", dn, len(sr.Entries))
			return sr.Entries[0], nil
		}
		debug("search %s with filter %s, not found", dn, filter)
		return nil, ErrNotFound
	}

	debug("ldap search %q, ERR: %s", dn, err)
	if le, ok := err.(*ldap.Error); ok && le.ResultCode == ldap.LDAPResultNoSuchObject {
		return nil, ErrNotFound
	}
	log.Printf("LDAP Search '%s' Error: %s", dn, err)
	return nil, err
}

func (ls *ldapSource) Authenticate(uid, passwd string) (err error) {
	err = ls.Bind(ls.UDN(uid), passwd)
	if err == ErrLogin && Domain != "" {
		upn := uid + "@" + Domain
		err = ls.Bind(upn, passwd)
	}
	return
}

func (ls *ldapSource) Bind(dn, passwd string) error {
	err := ls.opWithDN(dn, passwd, func(c ldap.Client) error {
		return nil
	})
	if err != nil {
		log.Printf("LDAP Bind failed for %s, reason: %s", dn, err)
		if le, ok := err.(*ldap.Error); ok {
			if le.ResultCode == ldap.LDAPResultInvalidCredentials ||
				le.ResultCode == ldap.LDAPResultInvalidDNSyntax {
				return ErrLogin
			}
		}
		return err
	}

	debug("bind(%s, ***) ok", dn)
	return nil
}

type opFunc func(c ldap.Client) error

// opWithMan admin operate
func (ls *ldapSource) opWithMan(op opFunc) error {
	return ls.opWithDN(ls.BindDN, ls.Passwd, op)
}

func (ls *ldapSource) opWithDN(dn, passwd string, op opFunc) error {
	c, err := ls.cp.Get()
	if err == nil {
		defer ls.cp.Put(c)
		err = c.Bind(dn, passwd)
		if err == nil {
			debug("conn from %s (len %d, idle %d) and bind(%s) ok", ls.Addr, ls.cp.Len(), ls.cp.IdleLen(), dn)
			return op(c)
		}
		log.Printf("LDAP bind(%s) ERR %s", dn, err)
		return err
	}

	log.Printf("get LDAP client from pool error, %s:%v", ls.Addr, err)
	return err
}

func (ls *ldapSource) Group(cn string) (*ldap.Entry, error) {
	if isADsource {
		return ls.Entry(etADgroup.DN(cn), etADgroup.Filter, etADgroup.Attributes...)
	}
	return ls.Entry(etGroup.DN(cn), etGroup.Filter, etGroup.Attributes...)
}

func (ls *ldapSource) People(uid string) (*ldap.Entry, error) {
	if isADsource {
		return ls.Entry(etADuser.DN(uid), etADuser.Filter, etADuser.Attributes...)
	}
	return ls.Entry(etPeople.DN(uid), etPeople.Filter, etPeople.Attributes...)
}

// Entry return a special entry with dn and filter
func (ls *ldapSource) Entry(dn, filter string, attrs ...string) (*ldap.Entry, error) {
	var entry *ldap.Entry
	err := ls.opWithMan(func(c ldap.Client) (err error) {
		entry, err = ldapEntryGet(c, dn, filter, attrs...)
		return
	})
	return entry, err
}

// GetStaff : search an LDAP source if an entry (with uid) is valide and in the specific filter
func (ls *ldapSource) GetStaff(uid string) (staff *models.Staff, err error) {
	var entry *ldap.Entry
	entry, err = ls.People(uid)
	if err != nil {
		log.Printf("GetStaff(%s) ERR %s", uid, err)
		return nil, err
	}

	return entryToUser(entry), nil
}

func (ls *ldapSource) GetByDN(dn string) (staff *models.Staff, err error) {
	var et *entryType
	if isADsource {
		et = etADuser
	} else {
		et = etPeople
	}
	var entry *ldap.Entry
	entry, err = ls.Entry(dn, et.Filter, et.Attributes...)
	if err != nil {
		return
	}
	return entryToUser(entry), nil
}

func (ls *ldapSource) ListPaged(limit int) (staffs models.Staffs) {
	if limit < 1 {
		limit = 1
	}
	var et *entryType
	if isADsource {
		et = etADuser
	} else {
		et = etPeople
	}
	search := ldap.NewSearchRequest(
		ls.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		et.Filter,
		et.Attributes,
		nil)

	var (
		sr *ldap.SearchResult
	)
	err := ls.opWithMan(func(c ldap.Client) (err error) {
		sr, err = c.SearchWithPaging(search, uint32(limit))
		return
	})
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
	u = &models.Staff{
		DN:           entry.DN,
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
	if str := entry.GetAttributeValue("sAMAccountName"); str != "" && u.Uid == "" {
		u.Uid = str
	}
	if str := entry.GetAttributeValue("userPrincipalName"); str != "" && u.Email == "" {
		u.Email = str
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
		u.Created, err = time.Parse(TimeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	} else if str := entry.GetAttributeValue("createTimestamp"); str != "" {
		u.Created, err = time.Parse(TimeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	}
	if str := entry.GetAttributeValue("modifiedTime"); str != "" {
		u.Updated, err = time.Parse(TimeLayout, str)
		if err != nil {
			log.Printf("invalid time %s, ERR %s", str, err)
		}
	} else if str := entry.GetAttributeValue("modifyTimestamp"); str != "" {
		u.Updated, err = time.Parse(TimeLayout, str)
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

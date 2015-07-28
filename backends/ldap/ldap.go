package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
	"strings"
	"tuluu.com/liut/staffio/models"
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
}

//Global LDAP directory pool
var (
	userDnFmt         = "uid=%s,ou=people,%s"
	defaultAttributes = []string{"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "employeeNumber", "employeeType", "description", "title"}
	defaultFilter     = "(objectclass=inetOrgPerson)"
	AuthenSource      []*LdapSource
)

// Add a new source (LDAP directory) to the global pool
func AddSource(addr, base string) *LdapSource {
	if base == "" {
		log.Fatal("ldap base is empty")
	}
	ls := &LdapSource{
		Addr:       addr,
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
		ls.c, err = ldap.DialTLS("tcp", ls.Addr, nil)
	} else {
		ls.c, err = ldap.Dial("tcp", ls.Addr)
	}

	if err != nil {
		log.Printf("LDAP Connect error, %s:%v", ls.Addr, err)
		ls.Enabled = false
		return nil, err
	}

	if strings.HasPrefix(ls.Addr, "localhost:") {
		ls.c.Debug = true
	}

	return ls.c, nil
}

func (ls *LdapSource) Close() {
	if ls.c != nil {
		ls.c.Close()
		ls.c = nil
	}
}

func Login(name, passwd string) (r bool, staff *models.Staff) {
	r = false
	for _, ls := range AuthenSource {
		r, staff = ls.SearchEntry(name, passwd)
		if r {
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

// searchEntry : search an LDAP source if an entry (name, passwd) is valide and in the specific filter
func (ls *LdapSource) SearchEntry(name, passwd string) (bool, *models.Staff) {
	l, err := ls.dial()
	if err != nil {
		return false, nil
	}

	dn := fmt.Sprintf(userDnFmt, name, ls.Base)
	err = l.Bind(dn, passwd)
	if err != nil {
		log.Printf("LDAP Authan failed for %s, reason: %s", dn, err.Error())
		return false, nil
	}

	search := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		ls.Filter,
		ls.Attributes,
		nil)
	sr, err := l.Search(search)
	if err != nil {
		log.Printf("LDAP Authen OK but not in filter %s", name)
		return false, nil
	}
	log.Printf("LDAP Authen OK: %s", name)
	if len(sr.Entries) > 0 {
		return true, entryToUser(sr.Entries[0])
		// cn := sr.Entries[0].GetAttributeValue(ls.AttributeUsername)
		// name := sr.Entries[0].GetAttributeValue(ls.AttributeName)
		// sn := sr.Entries[0].GetAttributeValue(ls.AttributeSurname)
		// mail := sr.Entries[0].GetAttributeValue(ls.AttributeMail)
		// return cn, name, sn, mail, true
	}
	return true, nil
}

func (ls *LdapSource) ListPaged(limit int) (staffs []*models.Staff) {
	l, err := ls.dial()
	if err != nil {
		return nil
	}

	err = l.Bind(ls.BindDN, ls.Passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
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

	sr, err := l.SearchWithPaging(search, uint32(limit))
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
	u.SurName = entry.GetAttributeValue("sn")
	u.GivenName = entry.GetAttributeValue("givenName")
	u.CommonName = entry.GetAttributeValue("cn")
	u.Email = entry.GetAttributeValue("mail")
	u.DisplayName = entry.GetAttributeValue("displayName")
	u.Mobile = entry.GetAttributeValue("mobile")
	u.EmployeeNumber = entry.GetAttributeValue("employeeNumber")
	u.Description = entry.GetAttributeValue("description")
	return
}

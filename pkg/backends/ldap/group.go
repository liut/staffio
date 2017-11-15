package ldap

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-ldap/ldap"

	"github.com/liut/staffio/pkg/models"
)

var (
	groupSuffix = "ou=groups"
	groupDnFmt  = "cn=%s,%s,%s"
	groupLimit  = 20
)

func (s *LDAPStore) AllGroup() (data []models.Group) {
	var err error
	for _, ls := range s.sources {
		data, err = ls.SearchGroup("")
		if err == nil {
			return
		}
	}
	if err == nil {
		err = ErrNotFound
	}
	return
}

func (s *LDAPStore) GetGroup(name string) (group *models.Group, err error) {
	// log.Printf("Search group %s", name)
	for _, ls := range s.sources {
		var data []models.Group
		data, err = ls.SearchGroup(name)
		if err == nil {
			group = &data[0]
			return
		}
		log.Printf("search group %q from %s error: %s", name, ls.Addr, err)
	}
	log.Printf("group %s not found", name)
	if err == nil {
		err = ErrNotFound
	}
	return
}

func (ls *ldapSource) GDN(name string) string {
	return fmt.Sprintf(groupDnFmt, name, groupSuffix, ls.Base)
}

func (ls *ldapSource) SearchGroup(name string) (data []models.Group, err error) {
	l, err := ls.dial()
	if err != nil {
		return nil, err
	}

	err = l.Bind(ls.BindDN, ls.Passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		return nil, err
	}

	var (
		dn string
	)
	if name == "" { // all
		dn = fmt.Sprintf("%s,%s", groupSuffix, ls.Base)
	} else {
		dn = ls.GDN(name)
	}

	search := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		etGroup.Filter,
		etGroup.Attributes,
		nil)
	sr, err := ls.c.SearchWithPaging(search, uint32(groupLimit))
	if err != nil {
		log.Printf("LDAP search group error: %s", err)
		return nil, err
	}

	if len(sr.Entries) > 0 {
		data = make([]models.Group, len(sr.Entries))
		for i, entry := range sr.Entries {
			data[i] = entryToGroup(entry)
		}
	}

	return
}

func entryToGroup(entry *ldap.Entry) (g models.Group) {
	for _, attr := range entry.Attributes {
		if attr.Name == "cn" {
			g.Name = attr.Values[0]
		} else if attr.Name == "member" {
			g.Members = make([]string, len(attr.Values))
			for j, _dn := range attr.Values {
				g.Members[j] = _dn[strings.Index(_dn, "=")+1 : strings.Index(_dn, ",")]
			}
		}
	}
	return
}

func (s *LDAPStore) SaveGroup(group *models.Group) error {
	// TODO:
	return nil
}

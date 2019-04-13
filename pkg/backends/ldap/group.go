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

func (s *LDAPStore) AllGroup() (data []models.Group, err error) {
	for _, ls := range s.sources {
		data, err = ls.SearchGroup("")
		if err == nil {
			return
		}
	}
	return
}

func (s *LDAPStore) GetGroup(name string) (group *models.Group, err error) {
	// log.Printf("Search group %s", name)
	for _, ls := range s.sources {
		var entry *ldap.Entry
		entry, err = ls.Group(name)
		if err == nil {
			group = entryToGroup(entry)
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
	return etGroup.DN(name)
	// return fmt.Sprintf(groupDnFmt, name, groupSuffix, ls.Base)
}

func (ls *ldapSource) SearchGroup(name string) (data []models.Group, err error) {
	var (
		dn string
	)
	if name == "" { // all
		dn = fmt.Sprintf("%s,%s", groupSuffix, ls.Base)
	} else {
		dn = ls.GDN(name)
	}

	var sr *ldap.SearchResult
	err = ls.opWithMan(func(c ldap.Client) (err error) {
		search := ldap.NewSearchRequest(
			dn,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			etGroup.Filter,
			etGroup.Attributes,
			nil)
		sr, err = c.SearchWithPaging(search, uint32(groupLimit))
		return
	})

	if err != nil {
		log.Printf("LDAP search group error: %s", err)
		return
	}

	if len(sr.Entries) > 0 {
		data = make([]models.Group, len(sr.Entries))
		for i, entry := range sr.Entries {
			g := entryToGroup(entry)
			data[i] = *g
		}
	}

	return
}

func entryToGroup(entry *ldap.Entry) (g *models.Group) {
	g = new(models.Group)
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

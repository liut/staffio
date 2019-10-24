package ldap

import (
	"log"
	"strings"

	"github.com/go-ldap/ldap"
)

const (
	groupAdminDefault = "keeper"
	groupAdminAD      = "Administrators"
)

var (
	groupLimit = 20
)

func (s *LDAPStore) AllGroup() (data []Group, err error) {
	for _, ls := range s.sources {
		data, err = ls.SearchGroup("")
		if err == nil {
			return
		}
	}
	return
}

func (s *LDAPStore) GetGroup(name string) (group *Group, err error) {
	// debug("Search group %s", name)
	for _, ls := range s.sources {
		var entry *ldap.Entry
		entry, err = ls.getGroupEntry(name)
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

func (ls *ldapSource) SearchGroup(name string) (data []Group, err error) {
	var (
		dn string
	)
	var et *entryType
	if ls.isAD {
		et = etADgroup
	} else {
		et = etGroup
	}
	if name == "" { // all
		dn = ls.Base
	} else {
		dn = et.DN(name, ls.Base)
	}

	var sr *ldap.SearchResult
	err = ls.opWithMan(func(c ldap.Client) (err error) {
		search := ldap.NewSearchRequest(
			dn,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			et.Filter,
			et.Attributes,
			nil)
		sr, err = c.SearchWithPaging(search, uint32(groupLimit))
		return
	})

	if err != nil {
		log.Printf("LDAP search group error: %s", err)
		return
	}

	if len(sr.Entries) > 0 {
		data = make([]Group, len(sr.Entries))
		for i, entry := range sr.Entries {
			g := entryToGroup(entry)
			data[i] = *g
		}
	}

	return
}

func entryToGroup(entry *ldap.Entry) (g *Group) {
	g = new(Group)
	for _, attr := range entry.Attributes {
		if attr.Name == "cn" || attr.Name == "name" {
			g.Name = attr.Values[0]
		} else if attr.Name == "member" {
			g.Members = make([]string, len(attr.Values))
			for j, _dn := range attr.Values {
				g.Members[j] = _dn[strings.Index(_dn, "=")+1 : strings.Index(_dn, ",")]
			}
		}
	}
	// debug("group %q", g)
	return
}

func (s *LDAPStore) SaveGroup(group *Group) error {
	for _, ls := range s.sources {
		err := ls.saveGroup(group)
		if err != nil {
			log.Printf("save group %q ERR %s", group.Name, err)
			return err
		}
	}
	return nil
}

func (ls *ldapSource) saveGroup(group *Group) error {
	if ls.isAD {
		return ErrUnsupport
	}
	err := ls.opWithMan(func(c ldap.Client) error {
		gdn := etGroup.DN(group.Name, ls.Base)
		var members []string
		for _, m := range group.Members {
			members = append(members, ls.UDN(m))
		}
		_, err := ldapEntryGet(c, gdn, etGroup.Filter, etGroup.Attributes...)
		if err == nil { // update
			mr := ldap.NewModifyRequest(gdn, nil)
			mr.Replace("member", members)
			debug("change group %v", mr)
			err = c.Modify(mr)
		}
		if err == ErrNotFound { // create
			ar := ldap.NewAddRequest(gdn, nil)
			etGroup.prepareTo(group.Name, ar)
			ar.Attribute("member", members)
			debug("add group %v", ar)
			err = c.Add(ar)
		}
		return err
	})
	return err
}

func (s *LDAPStore) EraseGroup(name string) error {
	for _, ls := range s.sources {
		err := ls.eraseGroup(name)
		if err != nil {
			log.Printf("save group %q ERR %s", name, err)
			return err
		}
	}
	return nil
}

func (ls *ldapSource) eraseGroup(name string) error {
	if ls.isAD {
		return ErrUnsupport
	}
	err := ls.opWithMan(func(c ldap.Client) error {
		dr := ldap.NewDelRequest(etGroup.DN(name, ls.Base), nil)
		return c.Del(dr)
	})
	return err
}

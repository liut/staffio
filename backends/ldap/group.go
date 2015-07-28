package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
	"strings"
	"tuluu.com/liut/staffio/models"
)

var (
	groupDnFmt = "cn=%s,ou=groups,%s"
)

func SearchGroup(name string) *models.Group {
	for _, ls := range AuthenSource {
		group, err := ls.SearchGroup(name)
		if err == nil {
			return group
		}
		log.Printf("search group from %s error: %s", ls.Addr, err)
	}
	return models.EmptyGroup
}

func (ls *LdapSource) SearchGroup(name string) (*models.Group, error) {
	l, err := ls.dial()
	if err != nil {
		return nil, err
	}
	defer l.Close()

	err = l.Bind(ls.BindDN, ls.Passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		return nil, err
	}

	search := ldap.NewSearchRequest(
		fmt.Sprintf(groupDnFmt, name, ls.Base),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectclass=groupOfNames)",
		[]string{"member"},
		nil)
	sr, err := l.Search(search)
	if err != nil {
		log.Printf("LDAP search error: %s", err)
		return nil, err
	}

	vals := sr.Entries[0].GetAttributeValues("member")

	members := make([]string, len(vals))
	for i, dn := range vals {
		members[i] = dn[strings.Index(dn, "="):strings.Index(dn, ",")]
	}

	return &models.Group{name, members}, nil
}

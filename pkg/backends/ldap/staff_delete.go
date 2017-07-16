package ldap

import (
	"log"

	"github.com/go-ldap/ldap"
)

func (s *LDAPStore) Delete(uid string) (err error) {
	for _, ls := range s.sources {
		err = ls.DeleteStaff(uid)
		if err != nil {
			return
		}
	}
	return
}

func (ls *ldapSource) DeleteStaff(uid string) (err error) {
	err = ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		return
	}
	dn := ls.UDN(uid)
	delRequest := ldap.NewDelRequest(dn, nil)
	err = ls.c.Del(delRequest)

	if err != nil {
		log.Printf("DeleteStaff %q Err: %s", uid, err)
	}

	return
}

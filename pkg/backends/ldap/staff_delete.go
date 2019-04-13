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
	if err = ls.Delete(ls.UDN(uid)); err != nil {
		log.Printf("DeleteStaff %q Err: %s", uid, err)
	}

	return
}

func (ls *ldapSource) Delete(dn string) error {
	return ls.opWithMan(func(c ldap.Client) (err error) {
		err = ldapEntryDel(c, dn)
		if err != nil {
			log.Printf("LDAP delete(%s) ERR %s", dn, err)
		}
		debug("delete %q, err %v", dn, err)
		return
	})
}

func ldapEntryDel(c ldap.Client, dn string) error {
	delRequest := ldap.NewDelRequest(dn, nil)
	return c.Del(delRequest)
}

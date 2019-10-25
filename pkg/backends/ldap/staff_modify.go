package ldap

import (
	// "fmt"
	"log"

	"github.com/go-ldap/ldap"
)

// ModifyBySelf ...
func (s *Store) ModifyBySelf(uid, password string, staff *People) (err error) {
	for _, ls := range s.sources {
		err = ls.Modify(uid, password, staff)
		if err != nil {
			log.Printf("Modify at %s ERR: %s", ls.Addr, err)
		}
	}
	return
}

func (ls *ldapSource) Modify(uid, password string, staff *People) error {

	debug("modify self %s staff: %v", uid, staff)

	userdn := ls.UDN(uid)
	return ls.opWithDN(userdn, password, func(c ldap.Client) (err error) {
		entry, err := ldapEntryGet(c, userdn, etPeople.Filter, etPeople.Attributes...)
		if err != nil {
			return err
		}

		modify := makeModifyRequest(userdn, entry, staff)

		if err = c.Modify(modify); err != nil {
			log.Printf("Modify ERROR: %s\n", err)
		}
		debug("modified %q, err %v", userdn, err)
		return nil
	})

}

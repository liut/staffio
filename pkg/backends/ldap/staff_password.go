package ldap

import (
	"log"

	"github.com/go-ldap/ldap"
)

func (s *LDAPStore) PasswordChange(uid, oldPasswd, newPasswd string) (err error) {
	for _, ls := range s.sources {
		err = ls.PasswordChange(uid, oldPasswd, newPasswd)
		if err != nil {
			log.Printf("PasswordChange at %s ERR: %s", ls.Addr, err)
			break
		}
	}
	return
}

func (ls *ldapSource) PasswordChange(uid, oldPasswd, newPasswd string) error {
	userdn := ls.UDN(uid)
	c, err := ls.cp.Get()
	defer ls.cp.Put(c)
	if err == nil {
		pmr := ldap.NewPasswordModifyRequest(userdn, oldPasswd, newPasswd)
		_, err := c.PasswordModify(pmr)
		if err != nil {
			log.Printf("PasswordModify(%s) ERR: %s", uid, err)
			return err
		}
		debug("PasswordModify(%s) OK", uid)
	}

	return err
}

func (s *LDAPStore) PasswordReset(uid, passwd string) (err error) {
	for _, ls := range s.sources {
		err = ls.PasswordReset(uid, passwd)
		if err != nil {
			log.Printf("PasswordReset at %s ERR: %s", ls.Addr, err)
			break
		}
	}
	return
}

// password reset by administrator
func (ls *ldapSource) PasswordReset(uid, newPasswd string) error {
	dn := ls.UDN(uid)
	return ls.opWithMan(func(c ldap.Client) error {
		passwordModifyRequest := ldap.NewPasswordModifyRequest(dn, "", newPasswd)
		_, err := c.PasswordModify(passwordModifyRequest)
		if err != nil {
			log.Printf("PasswordModify(%s) ERR: %s", uid, err)
			return err
		}
		debug("PasswordModify(%s) OK", uid)
		return nil
	})
}

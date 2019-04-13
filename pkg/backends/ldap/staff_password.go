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
		}
	}
	return
}

func (ls *ldapSource) PasswordChange(uid, oldPasswd, newPasswd string) error {
	userdn := ls.UDN(uid)
	return ls.opWithDN(userdn, oldPasswd, func(c ldap.Client) error {
		passwordModifyRequest := ldap.NewPasswordModifyRequest(userdn, oldPasswd, newPasswd)
		passwordModifyResponse, err := c.PasswordModify(passwordModifyRequest)

		if err != nil {
			log.Printf("PasswordModify ERR: %s", err)
			return err
		}
		log.Printf("passwordModifyResponse: %v", passwordModifyResponse)
		return nil
	})
}

func (s *LDAPStore) PasswordReset(uid, passwd string) (err error) {
	for _, ls := range s.sources {
		err = ls.PasswordReset(uid, passwd)
		if err != nil {
			log.Printf("PasswordReset at %s ERR: %s", ls.Addr, err)
		}
	}
	return
}

// password reset by administrator
func (ls *ldapSource) PasswordReset(uid, newPasswd string) error {
	err := ls.opWithMan(func(c ldap.Client) error {
		dn := ls.UDN(uid)
		passwordModifyRequest := ldap.NewPasswordModifyRequest(dn, "", newPasswd)
		passwordModifyResponse, err := c.PasswordModify(passwordModifyRequest)
		log.Printf("passwordModifyResponse: %v", passwordModifyResponse)
		return err
	})

	if err != nil {
		log.Printf("PasswordModify ERR: %s", err)
		return err
	}

	return nil
}

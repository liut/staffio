package ldap

import (
	"github.com/go-ldap/ldap"
	"log"
)

func PasswordChange(uid, oldPasswd, newPasswd string) (err error) {
	for _, ls := range AuthenSource {
		err = ls.PasswordChange(uid, oldPasswd, newPasswd)
		if err != nil {
			log.Printf("PasswordChange at %s ERR: %s", ls.Addr, err)
		}
	}
	return
}

func (ls *LdapSource) PasswordChange(uid, oldPasswd, newPasswd string) error {
	userdn := ls.UDN(uid)
	err := ls.Bind(userdn, oldPasswd, true)
	if err != nil {
		return err
	}
	passwordModifyRequest := ldap.NewPasswordModifyRequest(userdn, oldPasswd, newPasswd)
	passwordModifyResponse, err := ls.c.PasswordModify(passwordModifyRequest)

	if err != nil {
		log.Printf("PasswordModify ERR: %s", err)
		return err
	}

	log.Printf("passwordModifyResponse: %v", passwordModifyResponse)
	return nil
}

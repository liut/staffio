package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"log"
)

func Modify(uid, password string, values map[string]string) (err error) {
	for _, ls := range AuthenSource {
		err = ls.Modify(uid, password, values)
		if err != nil {
			log.Printf("Modify at %s ERR: %s", ls.Addr, err)
		}
	}
	return
}

func (ls *LdapSource) Modify(uid, password string, values map[string]string) error {
	if ls.Debug {
		log.Printf("change profile for %s values: %v", uid, values)
	}
	userdn := ls.UDN(uid)
	err := ls.Bind(userdn, password, true)
	if err != nil {
		return ErrLogin
	}
	entry, err := ls.getEntry(userdn)
	if err != nil {
		return err
	}

	modify := ldap.NewModifyRequest(entry.DN)
	changed := make(map[string]bool)
	for k, v := range values {
		if v == "" {
			continue
		}
		vals := entry.GetAttributeValues(k)
		if len(vals) == 0 {
			changed[k] = true
			modify.Add(k, []string{v})
		} else {
			if vals[0] != v {
				changed[k] = true
				modify.Replace(k, []string{v})
			}
		}
	}

	if len(changed) == 0 {
		if ls.Debug {
			log.Printf("nothing changed for %s", uid)
		}
		return nil
	}

	_, sok := changed["sn"]
	_, gok := changed["givenName"]
	if sok && gok {
		modify.Replace("cn", []string{fmt.Sprintf("%s%s", values["sn"], values["givenName"])})
	}

	if err := ls.c.Modify(modify); err != nil {
		log.Printf("Modify ERROR: %s\n", err.Error())
	}

	return nil
}

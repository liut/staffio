package ldap

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
	"github.com/jsimonetti/pwscheme/ssha"
	garbler "github.com/michaelbironneau/garbler/lib"
	"github.com/wealthworks/csmtp"

	"lcgc/platform/staffio/models"
)

var (
	// simpleSecurityObject
	objectClassPeople = []string{"top", "simpleSecurityObject", "PersonExt", "uidObject", "inetOrgPerson"}
)

func StoreStaff(staff *models.Staff) (err error) {
	var (
		pass, hash string
	)
	if staff.Passwd == "" {
		reqs := garbler.PasswordStrengthRequirements{MinimumTotalLength: 16, Digits: 10}
		pass, err = garbler.NewPassword(&reqs)
		if err != nil {
			log.Printf("garbler.NewPassword err %s", err)
			return
		}
		hash, err = ssha.Generate(pass, 4)
		if err != nil {
			log.Printf("ssha.Generate err %s", err)
			return
		}
		log.Printf("gen new passwd %q", pass)
		staff.Passwd = hash
	}
	var isNew bool
	for _, ls := range AuthenSource {
		isNew, err = ls.StoreStaff(staff)
		if err != nil {
			log.Printf("StoreStaff at %s ERR: %s", ls.Addr, err)
			return
		}
	}
	if isNew {
		message := fmt.Sprintf("Your new password is <strong>%s</strong>.", hash)
		csmtp.SendMail("Welcome!", message, staff.Email)
	}
	return
}

func (ls *LdapSource) StoreStaff(staff *models.Staff) (isNew bool, err error) {
	uid := staff.Uid
	err = ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		return
	}
	dn := ls.UDN(uid)
	_, err = ls.getEntry(ls.UDN(uid))
	if err == nil {
		// TODO: update
		mr := ldap.NewModifyRequest(dn)
		mr.Replace("sn", []string{staff.Surname})
		mr.Replace("givenName", []string{staff.GivenName})
		mr.Replace("cn", []string{staff.CommonName})
		mr.Replace("mail", []string{staff.Email})

		if staff.Mobile != "" {
			mr.Replace("mobile", []string{staff.Mobile})
		}
		if staff.EmployeeNumber != "" {
			mr.Replace("employeeNumber", []string{staff.EmployeeNumber})
		}

		if staff.Description != "" {
			mr.Replace("description", []string{staff.Description})
		}

		err = ls.c.Modify(mr)
		if err != nil {
			log.Printf("add err %s", err)
		}
		return
	}
	if err == ErrNotFound {
		isNew = true
		ar := ldap.NewAddRequest(dn)
		ar.Attribute("objectClass", objectClassPeople)
		ar.Attribute("uid", []string{uid})
		ar.Attribute("sn", []string{staff.Surname})
		ar.Attribute("givenName", []string{staff.GivenName})
		ar.Attribute("cn", []string{staff.CommonName})
		ar.Attribute("mail", []string{staff.Email})
		if staff.Mobile != "" {
			ar.Attribute("mobile", []string{staff.Mobile})
		}

		if staff.EmployeeNumber != "" {
			ar.Attribute("employeeNumber", []string{staff.EmployeeNumber})
		}
		if staff.EmployeeType != "" {
			ar.Attribute("employeeType", []string{staff.EmployeeType})
		}
		if staff.Description != "" {
			ar.Attribute("description", []string{staff.Description})
		}

		ar.Attribute("userPassword", []string{staff.Passwd})

		err = ls.c.Add(ar)
		if err != nil {
			log.Printf("add err %s", err)
		}
		return
	}

	log.Printf("getEntry err %s", err)

	return
}

/*
uid
sn
givenName
cn
mail
displayName
mobile
employeeNumber
employeeType
description
*/

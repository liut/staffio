package ldap

import (
	"log"
	"strconv"
	"time"

	"github.com/go-ldap/ldap"

	"github.com/liut/staffio/pkg/models"
)

func (ls *ldapSource) storeStaff(staff *models.Staff) (isNew bool, err error) {
	err = ls.opWithMan(func(c ldap.Client) (err error) {
		dn := ls.UDN(staff.Uid)
		var entry *ldap.Entry
		entry, err = ldapEntryGet(c, ls.UDN(staff.Uid), etPeople.Filter, etPeople.Attributes...)
		if err == nil {
			// :update
			mr := makeModifyRequest(dn, entry, staff)
			eidStr := strconv.Itoa(staff.EmployeeNumber)
			if staff.EmployeeNumber > 0 && eidStr != entry.GetAttributeValue("employeeNumber") {
				mr.Replace("employeeNumber", []string{eidStr})
			}
			if staff.EmployeeType != entry.GetAttributeValue("employeeType") {
				mr.Replace("employeeType", []string{staff.EmployeeType})
			}
			err = c.Modify(mr)
			if err != nil {
				log.Printf("modify %v ERR %s", mr, err)
			}
			return
		}
		if err == ErrNotFound {
			isNew = true
			ar := makeAddRequest(dn, staff)
			err = c.Add(ar)
			if err != nil {
				log.Printf("add %v ERR %s", ar, err)
			}
			return
		}
		log.Printf("storeStaff %s ERR %s", staff.Uid, err)

		return
	})

	return
}

func makeAddRequest(dn string, staff *models.Staff) *ldap.AddRequest {
	ar := ldap.NewAddRequest(dn, nil)
	ar.Attribute("objectClass", objectClassPeople)
	ar.Attribute("uid", []string{staff.Uid})
	ar.Attribute("cn", []string{staff.GetCommonName()})
	if staff.Surname != "" {
		ar.Attribute("sn", []string{staff.Surname})
	}
	if staff.GivenName != "" {
		ar.Attribute("givenName", []string{staff.GivenName})
	}

	if staff.Email != "" {
		ar.Attribute("mail", []string{staff.Email})
	}

	if staff.Nickname != "" {
		ar.Attribute("displayName", []string{staff.Nickname})
	}
	if staff.Mobile != "" {
		ar.Attribute("mobile", []string{staff.Mobile})
	}

	if staff.EmployeeNumber > 0 {
		ar.Attribute("employeeNumber", []string{strconv.Itoa(staff.EmployeeNumber)})
	}
	if staff.EmployeeType != "" {
		ar.Attribute("employeeType", []string{staff.EmployeeType})
	}
	if staff.Gender != models.Unknown {
		ar.Attribute("gender", []string{staff.Gender.String()[0:1]})
	}
	if staff.Birthday != "" {
		ar.Attribute("dateOfBirth", []string{staff.Birthday})
	}
	if staff.Description != "" {
		ar.Attribute("description", []string{staff.Description})
	}
	if staff.AvatarPath != "" {
		ar.Attribute("avatarPath", []string{staff.AvatarPath})
	}
	if staff.JoinDate != "" {
		ar.Attribute("dateOfJoin", []string{staff.JoinDate})
	}

	// if staff.Passwd != "" {
	// 	ar.Attribute("userPassword", []string{staff.Passwd})
	// }

	return ar
}

func makeModifyRequest(dn string, entry *ldap.Entry, staff *models.Staff) *ldap.ModifyRequest {
	mr := ldap.NewModifyRequest(dn, nil)
	mr.Replace("objectClass", objectClassPeople)
	if staff.Surname != entry.GetAttributeValue("sn") {
		mr.Replace("sn", []string{staff.Surname})
	}
	if staff.GivenName != entry.GetAttributeValue("givenName") {
		mr.Replace("givenName", []string{staff.GivenName})
	}
	if staff.CommonName != entry.GetAttributeValue("cn") {
		mr.Replace("cn", []string{staff.GetCommonName()})
	}
	if len(staff.Nickname) > 0 && staff.Nickname != entry.GetAttributeValue("displayName") {
		mr.Replace("displayName", []string{staff.Nickname})
	}
	if len(staff.Email) > 0 && staff.Email != entry.GetAttributeValue("mail") {
		mr.Replace("mail", []string{staff.Email})
	}
	if len(staff.Mobile) > 0 && staff.Mobile != entry.GetAttributeValue("mobile") {
		mr.Replace("mobile", []string{staff.Mobile})
	}
	if len(staff.AvatarPath) > 0 && staff.AvatarPath != entry.GetAttributeValue("avatarPath") {
		mr.Replace("avatarPath", []string{staff.AvatarPath})
	}
	if staff.Gender != models.Unknown {
		mr.Replace("gender", []string{staff.Gender.String()[0:1]})
	}
	if len(staff.Birthday) > 0 && staff.Birthday != entry.GetAttributeValue("dateOfBirth") {
		mr.Replace("dateOfBirth", []string{staff.Birthday})
	}
	if len(staff.Description) > 0 && staff.Description != entry.GetAttributeValue("description") {
		mr.Replace("description", []string{staff.Description})
	}
	mr.Replace("modifiedTime", []string{time.Now().Format(TimeLayout)})
	return mr
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

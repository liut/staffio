package ldap

import (
	"fmt"
)

// LDAP config
type Config struct {
	Addr, Base   string
	Bind, Passwd string
	Filter       string
	Attributes   []string
}

type entryType struct {
	PK, OC     string
	Filter     string
	Attributes []string
}

func NewEentryType(pk, oc string, attrs ...string) (et *entryType) {
	et = &entryType{
		PK:         pk,
		OC:         oc,
		Attributes: attrs,
	}
	et.Filter = fmt.Sprintf("(objectclass=%s)", et.OC)

	return
}

func (et *entryType) DN(name string) string {
	switch et {
	case etGroup:
		return DN(et.PK, name, etParent.DN("groups"))
	case etPeople:
		return DN(et.PK, name, etParent.DN("people"))
	case etBase:
		return Base
	}
	// parent
	return DN(et.PK, name, Base)
}

func DN(pk, name, parent string) string {
	return fmt.Sprintf("%s=%s,%s", pk, name, parent)
}

var (
	Base   = "dc=mydomain,dc=net"
	Domain = "mydomain.net"

	etBase   = NewEentryType("dc", "dcObject", "dc", "o")
	etParent = NewEentryType("ou", "organizationalUnit", "ou")
	etGroup  = NewEentryType("cn", "groupOfNames", "cn", "member")
	etPeople = NewEentryType("uid", "inetOrgPerson",
		"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "description",
		"avatarPath", "dateOfBirth", "gender", "employeeNumber", "employeeType", "title")
)

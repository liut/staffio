package ldap

import (
	"fmt"
	"os"
)

// Config LDAP config
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

const (
	timeLayout = "20060102150405Z"
)

var (
	Base   string
	Domain string

	etBase   = NewEentryType("dc", "dcObject", "dc", "o")
	etParent = NewEentryType("ou", "organizationalUnit", "ou")
	etGroup  = NewEentryType("cn", "groupOfNames", "cn", "member")
	etPeople = NewEentryType("uid", "inetOrgPerson",
		"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "description",
		"createdTime", "modifiedTime", "createTimestamp", "modifyTimestamp", "jpegPhoto",
		"avatarPath", "dateOfBirth", "gender", "employeeNumber", "employeeType", "title")

	objectClassPeople = []string{"top", "staffioPerson", "uidObject", "inetOrgPerson"}

	PoolSize = 10
)

func init() {
	Base = envOr("LDAP_BASE_DN", "dc=mydomain,dc=net")
	Domain = envOr("LDAP_DOMAIN", "mydomain.net")
}

func envOr(key, dft string) string {
	v := os.Getenv(key)
	if v == "" {
		return dft
	}
	return v
}

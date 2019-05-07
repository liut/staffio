package ldap

import (
	"fmt"
	"os"
)

/*

	cfg := NewConfig()
	store, err := NewStore(cfg)
	if err != nil {
		log.Fatalf("new service ERR %s", err)
	}

*/

// Config LDAP config
type Config struct {
	Addr, Base   string
	Bind, Passwd string
	Domain       string
	PageSize     int
}

// NewConfig return default Config from Environment
func NewConfig() Config {
	return Config{
		Addr:     envOr("LDAP_HOSTS", envOr("STAFFIO_LDAP_HOSTS", "localhost")),
		Base:     envOr("LDAP_BASE_DN", envOr("STAFFIO_LDAP_BASE_DN", "dc=mydomain,dc=net")),
		Domain:   envOr("LDAP_DOMAIN", envOr("STAFFIO_EMAIL_DOMAIN", "mydomain.net")),
		Bind:     envOr("LDAP_BIND_DN", envOr("STAFFIO_LDAP_BIND_DN", "")),
		Passwd:   envOr("LDAP_PASSWD", envOr("STAFFIO_LDAP_PASS", "")),
		PageSize: DefaultPageSize,
	}
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
	if oc == "" {
		oc = "*"
	}
	et.Filter = fmt.Sprintf("(objectclass=%s)", oc)

	return
}

func (et *entryType) DN(name, base string) string {
	switch et {
	case etGroup:
		return DN(et.PK, name, etParent.DN("groups", base))
	case etPeople:
		return DN(et.PK, name, etParent.DN("people", base))
	case etADgroup:
		return "CN=" + name + ",CN=Builtin," + base
	case etADuser:
		return "CN=" + name + ",CN=Users," + base
	case etBase:
		return base
	}
	// parent
	return DN(et.PK, name, base)
}

type attributer interface {
	Attribute(attrType string, attrVals []string)
}

func (et *entryType) objectClasses() []string {
	if et.PK == "dc" {
		return []string{"dcObject", "organization", "top"}
	}
	return []string{et.OC, "top"}
}

func (et *entryType) prepareTo(name string, ar attributer) {
	ar.Attribute("objectClass", et.objectClasses())
	ar.Attribute(et.PK, []string{name})
	if et.PK == "dc" {
		ar.Attribute("o", []string{name})
	}
}

func DN(pk, name, parent string) string {
	return fmt.Sprintf("%s=%s,%s", pk, name, parent)
}

const (
	TimeLayout = "20060102150405Z"
	DateLayout = "20060102"

	DefaultPageSize = 100
	DefaultPoolSize = 10
)

var (
	etBase   = NewEentryType("dc", "", "dc", "o", "instanceType")
	etParent = NewEentryType("ou", "organizationalUnit", "ou")
	etGroup  = NewEentryType("cn", "groupOfNames", "cn", "member")
	etPeople = NewEentryType("uid", "inetOrgPerson",
		"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "description",
		"createdTime", "modifiedTime", "createTimestamp", "modifyTimestamp", "jpegPhoto",
		"avatarPath", "dateOfBirth", "gender", "employeeNumber", "employeeType", "title")

	etADgroup = NewEentryType("cn", "group", "cn", "member", "name", "description", "instanceType")
	etADuser  = NewEentryType("cn", "user", "name", "sAMAccountName", "userPrincipalName",
		"uid", "gn", "sn", "cn", "displayName", "mail", "mobile", "description",
		"employeeNumber", "employeeType", "title", "jpegPhoto", "logonCount")

	objectClassPeople = []string{"top", "staffioPerson", "uidObject", "inetOrgPerson"}
)

func envOr(key, dft string) string {
	v := os.Getenv(key)
	if v == "" {
		return dft
	}
	return v
}

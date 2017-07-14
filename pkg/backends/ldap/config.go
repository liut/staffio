package ldap

// LDAP config
type Config struct {
	Addr, Base   string
	Bind, Passwd string
	Filter       string
	Attributes   []string
}

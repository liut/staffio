package config

import (
	"flag"

	envcfg "github.com/wealthworks/envflagset"
)

type config struct {
	Name    string
	Version string

	EmailDomain string

	SMTP struct {
		Host           string
		Port           int
		SenderName     string
		SenderEmail    string
		SenderPassword string
	}

	LDAP struct {
		Host     string
		Base     string
		BindDN   string
		Password string `ini:"pass"`
		Filter   string
	} `ini:"ldap"`

	Session struct {
		Name   string
		Domain string
		Secret string
		MaxAge int // cookie maxAge
	} `ini:"sess"`

	HttpListen   string
	BaseURL      string
	ResUrl       string
	PwdSecret    string
	Backend      struct{ DSN string }
	SentryDSN    string
	Root         string
	Debug        bool
	UserLifetime int // secends, user online time

	TokenGen struct { // JWT config
		Method string // disuse
		Key    string
	}

	CacheSize     int
	CacheLifetime uint
}

var (
	fs       *flag.FlagSet
	Settings *config = &config{}
)

func init() {
	Settings.Name = "staffio"
	Settings.Version = buildVersion

	fs = envcfg.New(Settings.Name, Settings.Version)

	fs.StringVar(&Settings.LDAP.Host, "ldap-host", "ldap://localhost:389", "ldap hostname")
	fs.StringVar(&Settings.LDAP.Base, "ldap-base", "", "ldap base")
	fs.StringVar(&Settings.LDAP.BindDN, "ldap-bind-dn", "", "ldap bind dn")
	fs.StringVar(&Settings.LDAP.Password, "ldap-pass", "", "ldap bind password")
	fs.StringVar(&Settings.LDAP.Filter, "ldap-user-filter", "(objectclass=inetOrgPerson)", "ldap search filter")

	fs.StringVar(&Settings.HttpListen, "http-listen", "localhost:5000", "bind address and port")
	fs.StringVar(&Settings.BaseURL, "prefix", "http://localhost:5000", "url prefix for self")
	fs.StringVar(&Settings.PwdSecret, "password-secret", "very secret", "the secret of password reset")
	fs.StringVar(&Settings.Session.Name, "sess-name", "staff_sess", "session name")
	fs.StringVar(&Settings.Session.Domain, "sess-domain", "", "session domain")
	fs.StringVar(&Settings.Session.Secret, "sess-secret", "very-secret", "session secret")
	fs.IntVar(&Settings.Session.MaxAge, "sess-maxage", 86400*7, "session cookie life time (in seconds)")
	fs.IntVar(&Settings.UserLifetime, "user-life", 2500, "user online life time (in seconds)")
	fs.StringVar(&Settings.ResUrl, "res-url", "/static/", "static resource url")

	fs.StringVar(&Settings.Backend.DSN, "backend-dsn", "postgres://staffio@localhost/staffio?sslmode=disable", "database dsn string for backend")
	fs.StringVar(&Settings.SentryDSN, "sentry-dsn", "", "SENTRY_DSN")
	fs.StringVar(&Settings.Root, "root", "./", "app root directory")
	fs.StringVar(&Settings.TokenGen.Key, "tokengen-key", "", "HMAC key for token generater")

	fs.StringVar(&Settings.EmailDomain, "email-domain", "example.net", "default email domain")
	fs.StringVar(&Settings.SMTP.Host, "smtp-host", "", "")
	fs.IntVar(&Settings.SMTP.Port, "smtp-port", 465, "")
	fs.StringVar(&Settings.SMTP.SenderName, "smtp-sender-name", "Notification", "")
	fs.StringVar(&Settings.SMTP.SenderEmail, "smtp-sender-email", "", "")
	fs.StringVar(&Settings.SMTP.SenderPassword, "smtp-sender-password", "", "")
	fs.BoolVar(&Settings.Debug, "debug", false, "app in debug mode")

	fs.IntVar(&Settings.CacheSize, "cache-size", 512*1024, "cache size")
	fs.UintVar(&Settings.CacheLifetime, "cache-life", 60, "cache lifetime in seconds")

}

func (c *config) Parse() {
	envcfg.Parse()
}

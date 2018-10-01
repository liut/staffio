package settings

import (
	"flag"

	envcfg "github.com/wealthworks/envflagset"
)

var (
	EmailDomain string
	EmailCheck  bool

	SMTP struct {
		Enabled        bool
		Host           string
		Port           int
		SenderName     string
		SenderEmail    string
		SenderPassword string
	}

	LDAP struct {
		Hosts    string
		Base     string
		BindDN   string
		Password string
		Filter   string
	}

	Session struct {
		Name   string
		Domain string
		Secret string
		MaxAge int // cookie maxAge
	}

	HttpListen   string
	BaseURL      string
	PwdSecret    string
	Backend      struct{ DSN string }
	SentryDSN    string
	Root         string
	FS           string
	CommonFormat string
	Debug        bool
	UserLifetime int // secends, user online time

	TokenGen struct { // JWT config
		Method string // disuse
		Key    string
	}

	CacheSize     int
	CacheLifetime uint

	WechatCorpID        string
	WechatContactSecret string
	WechatPortalSecret  string
	WechatPortalAgentID int
)

var (
	fs *flag.FlagSet
)

func init() {
	fs = envcfg.New(NAME, buildVersion)

	fs.StringVar(&LDAP.Hosts, "ldap-hosts", "ldap://localhost:389", "ldap hostname")
	fs.StringVar(&LDAP.Base, "ldap-base", "", "ldap base")
	fs.StringVar(&LDAP.BindDN, "ldap-bind-dn", "", "ldap bind dn")
	fs.StringVar(&LDAP.Password, "ldap-pass", "", "ldap bind password")
	fs.StringVar(&LDAP.Filter, "ldap-user-filter", "(objectclass=inetOrgPerson)", "ldap search filter")

	fs.StringVar(&HttpListen, "http-listen", "localhost:5000", "bind address and port")
	fs.StringVar(&BaseURL, "baseurl", "http://localhost:5000", "url base for self host")
	fs.StringVar(&PwdSecret, "password-secret", "very secret", "the secret of password reset")
	fs.StringVar(&Session.Name, "sess-name", "staff_sess", "session name")
	fs.StringVar(&Session.Domain, "sess-domain", "", "session domain")
	fs.StringVar(&Session.Secret, "sess-secret", "very-secret", "session secret")
	fs.IntVar(&Session.MaxAge, "sess-maxage", 86400*7, "session cookie life time (in seconds)")
	fs.IntVar(&UserLifetime, "user-life", 2500, "user online life time (in seconds)")

	fs.StringVar(&Backend.DSN, "backend-dsn", "postgres://staffio@localhost/staffio?sslmode=disable", "database dsn string for backend")
	fs.StringVar(&SentryDSN, "sentry-dsn", "", "SENTRY_DSN")
	fs.StringVar(&Root, "root", "./", "app root directory")
	fs.StringVar(&FS, "fs", "bind", "file system [bind | local]")
	fs.StringVar(&CommonFormat, "cn-fmt", "{sn}{gn}", "common name format, sn=surname, gn=given name")
	fs.StringVar(&TokenGen.Key, "tokengen-key", "", "HMAC key for token generater")

	fs.StringVar(&EmailDomain, "email-domain", "example.net", "default email domain")
	fs.BoolVar(&EmailCheck, "email-check", false, "check email unseen")
	fs.BoolVar(&SMTP.Enabled, "smtp-enabled", true, "enable smtp")
	fs.StringVar(&SMTP.Host, "smtp-host", "", "")
	fs.IntVar(&SMTP.Port, "smtp-port", 465, "")
	fs.StringVar(&SMTP.SenderName, "smtp-sender-name", "Notification", "")
	fs.StringVar(&SMTP.SenderEmail, "smtp-sender-email", "", "")
	fs.StringVar(&SMTP.SenderPassword, "smtp-sender-password", "", "")
	fs.BoolVar(&Debug, "debug", false, "app in debug mode")

	fs.IntVar(&CacheSize, "cache-size", 512*1024, "cache size")
	fs.UintVar(&CacheLifetime, "cache-life", 60, "cache lifetime in seconds")

	fs.StringVar(&WechatCorpID, "wechat-corpid", "", "wechat corpId")
	fs.StringVar(&WechatContactSecret, "wechat-contact-secret", "", "wechat secret of contacts")
	fs.StringVar(&WechatPortalSecret, "wechat-portal-secret", "", "wechat secret of portal(oauth)")
	fs.IntVar(&WechatPortalAgentID, "wechat-portal-agentid", 0, "wechat agentId of portal(oauth)")

}

func Parse() {
	envcfg.Parse()
}

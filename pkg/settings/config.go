package settings

import (
	"github.com/kelseyhightower/envconfig"
)

// Config config from env
type Config struct {
	HTTPListen string `envconfig:"HTTP_LISTEN" default:"localhost:3030"`
	BaseURL    string `envconfig:"BASEURL" default:"http://localhost:3030"`
	PwdSecret  string `envconfig:"PASSWORD_SECRET"`
	BackendDSN string `envconfig:"BACKEND_DSN"`
	SentryDSN  string `envconfig:"SENTRY_DSN"`

	Root  string `default:"./"`
	Debug bool

	TokenGenKey string `envconfig:"tokengen_key"`

	EmailDomain string `envconfig:"EMAIL_DOMAIN"`
	EmailCheck  bool   `envconfig:"EMAIL_CHECK"`

	SMTPEnabled        bool   `envconfig:"SMTP_ENABLED"`
	SMTPHost           string `envconfig:"SMTP_HOST"`
	SMTPPort           int    `envconfig:"SMTP_PORT" default:"465"`
	SMTPSenderName     string `envconfig:"SMTP_SENDER_NAME"`
	SMTPSenderEmail    string `envconfig:"SMTP_SENDER_EMAIL"`
	SMTPSenderPassword string `envconfig:"SMTP_SENDER_PASSWORD"`

	// LDAPHosts    string `envconfig:"LDAP_HOSTS" default:"localhost"`
	// LDAPBase     string `envconfig:"LDAP_BASE"`
	// LDAPDomain   string `envconfig:"LDAP_DOMAIN"`
	// LDAPBindDN   string `envconfig:"LDAP_BIND_DN"`
	// LDAPPassword string `envconfig:"LDAP_PASSWD"`

	WechatCorpID        string `envconfig:"wechat_corpid"`
	WechatContactSecret string `envconfig:"wechat_contact_secret"`
	WechatPortalSecret  string `envconfig:"wechat_portal_secret"`
	WechatPortalAgentID int    `envconfig:"wechat_portal_agentid"`

	InDevelop bool   `envconfig:"-"`
	Version   string `envconfig:"-"`
}

var Current *Config

func init() {
	Current = new(Config)
	err := envconfig.Process(NAME, Current)
	if err != nil {
		panic(err)
	}
	Current.InDevelop = IsDevelop()
	Current.Version = buildVersion
}

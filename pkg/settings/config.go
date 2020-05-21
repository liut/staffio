package settings

import (
	"github.com/kelseyhightower/envconfig"
)

// Config config from env
type Config struct {
	HTTPListen string `envconfig:"HTTP_LISTEN" default:"localhost:3030"`
	BaseURL    string `envconfig:"BASEURL" default:"http://localhost:3030"`
	PwdSecret  string `envconfig:"PASSWORD_SECRET"`
	BackendDSN string `envconfig:"BACKEND_DSN" default:"postgres://staffio@localhost/staffio?sslmode=disable"`
	SentryDSN  string `envconfig:"SENTRY_DSN"`

	RedisAddrs    []string `envconfig:"REDIS_ADDRS" `         // host:port,host:port
	RedisDB       int      `envconfig:"REDIS_DB" default:"1"` // Redis DB 1
	RedisPassword string   `envconfig:"REDIS_PASSWROD"`

	Root  string `default:"./"`
	Debug bool   `default:"false"`

	TokenGenKey string `envconfig:"tokengen_key"`

	EmailDomain string `envconfig:"EMAIL_DOMAIN"`
	EmailCheck  bool   `envconfig:"EMAIL_CHECK" default:"false"`

	MailEnabled        bool   `envconfig:"SMTP_ENABLED" default:"false"`
	MailHost           string `envconfig:"SMTP_HOST"`
	MailPort           int    `envconfig:"SMTP_PORT" default:"465"`
	MailSenderName     string `envconfig:"SMTP_SENDER_NAME" default:"notify"`
	MailSenderEmail    string `envconfig:"SMTP_SENDER_EMAIL"`
	MailSenderPassword string `envconfig:"SMTP_SENDER_PASSWORD"`
	MailTLSEnabled     bool   `envconfig:"SMTP_TLS" default:"true"`

	LDAPHosts    string `envconfig:"LDAP_HOSTS" default:"localhost:389"`
	LDAPBase     string `envconfig:"LDAP_BASE"`
	LDAPDomain   string `envconfig:"LDAP_DOMAIN"` // used for AD
	LDAPBindDN   string `envconfig:"LDAP_BIND_DN"`
	LDAPPassword string `envconfig:"LDAP_PASSWD"`

	WechatCorpID        string `envconfig:"wechat_corpid"`
	WechatContactSecret string `envconfig:"wechat_contact_secret"`
	WechatPortalSecret  string `envconfig:"wechat_portal_secret"`
	WechatPortalAgentID int    `envconfig:"wechat_portal_agentid"`

	LarkAppID      string `envconfig:"lark_app_id"`
	LarkAppSecret  string `envconfig:"lark_app_secret"`
	LarkEncryptKey string `envconfig:"LARK_ENCRYPT_KEY"`
}

// Current ...
var Current *Config

func init() {
	Current = new(Config)
	err := envconfig.Process(NAME, Current)
	if err != nil {
		panic(err)
	}
}

// Usage print envs for config
func Usage() error {
	return envconfig.Usage(NAME, Current)
}

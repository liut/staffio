package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type config struct {
	Name    string
	Version string

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
	ResUrl       string
	Backend      struct{ DSN string }
	SentryDSN    string
	Root         string
	Debug        bool
	UserLifetime int // secends, user online time
}

var (
	_configFile  string
	_loaded      bool
	fs           *flag.FlagSet
	Settings     *config = &config{}
	printVersion bool
)

func init() {
	Settings.Name = "staffio"
	Settings.Version = buildVersion
	// Settings.HttpListen = "localhost:3000"
	// Settings.ResUrl = "/static/"
	// Settings.LDAP.Host = "localhost"

	fs = flag.NewFlagSet("staffio", flag.ExitOnError)

	fs.StringVar(&Settings.LDAP.Host, "ldap-host", "ldap://localhost:389", "ldap hostname")
	fs.StringVar(&Settings.LDAP.Base, "ldap-base", "", "ldap base")
	fs.StringVar(&Settings.LDAP.BindDN, "ldap-bind-dn", "", "ldap bind dn")
	fs.StringVar(&Settings.LDAP.Password, "ldap-pass", "", "ldap bind password")
	fs.StringVar(&Settings.LDAP.Filter, "ldap-user-filter", "(objectclass=inetOrgPerson)", "ldap search filter")
	fs.StringVar(&Settings.HttpListen, "http-listen", "localhost:5000", "bind address and port")
	fs.StringVar(&Settings.Session.Name, "sess-name", "staff_sess", "session name")
	fs.StringVar(&Settings.Session.Domain, "sess-domain", "", "session domain")
	fs.StringVar(&Settings.Session.Secret, "sess-secret", "very-secret", "session secret")
	fs.IntVar(&Settings.Session.MaxAge, "sess-maxage", 86400*30, "session cookie life time (in seconds)")
	fs.IntVar(&Settings.UserLifetime, "user-life", 2500, "user online life time (in seconds)")
	fs.StringVar(&Settings.ResUrl, "res-url", "/static/", "static resource url")
	fs.StringVar(&Settings.Backend.DSN, "backend-dsn", "postgres://staffio@localhost/staffio?sslmode=disable", "database dsn string for backend")
	fs.StringVar(&Settings.SentryDSN, "sentry-dsn", "", "SENTRY_DSN")
	fs.StringVar(&Settings.Root, "root", "./", "app root directory")
	fs.BoolVar(&Settings.Debug, "debug", false, "app in debug mode")
	fs.BoolVar(&printVersion, "version", false, "Print the version and exit")

}

func (c *config) Parse() {
	perr := fs.Parse(os.Args[1:])
	switch perr {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		os.Exit(2)
	}
	if len(fs.Args()) != 0 {
		log.Fatalf("'%s' is not a valid flag", fs.Arg(0))
	}

	if printVersion {
		fmt.Println("staffio version", Settings.Version)
		os.Exit(0)
	}

	err := SetFlagsFromEnv(fs, Settings.Name+"_")
	if err != nil {
		log.Fatalf("staffio: %v", err)
	}

}

// SetFlagsFromEnv parses all registered flags in the given flagset,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, have the prefix "PREFIX_", and any dashes are replaced by
// underscores - for example: some-flag => PREFIX_SOME_FLAG
func SetFlagsFromEnv(fs *flag.FlagSet, prefix string) error {
	var err error
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})

	if prefix == "" {
		prefix = "_"
	} else {
		prefix = strings.ToUpper(strings.Replace(prefix, "-", "_", -1))
	}

	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := prefix + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
			val := os.Getenv(key)
			if val != "" {
				if serr := fs.Set(f.Name, val); serr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
				}
			}
		}
	})
	return err
}

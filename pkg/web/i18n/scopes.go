package i18n

// Scope ...
type Scope string

// consts of scope
const (
	ScopeBasic   Scope = "basic"
	ScopeOpenID  Scope = "openid"
	ScopeProfile Scope = "profile"
	ScopeEmail   Scope = "email"
)

// Valid ...
func (s Scope) Valid() bool {
	switch s {
	case ScopeBasic, ScopeEmail, ScopeOpenID, ScopeProfile:
		return true
	default:
		return false
	}
}

func (s Scope) LabelP(p *Printer) string {
	switch s {
	case ScopeBasic:
		return p.Sprintf("Basic Information")
	case ScopeEmail:
		return p.Sprintf("Email address")
	case ScopeOpenID:
		return p.Sprintf("OpenID Connect")
	case ScopeProfile:
		return p.Sprintf("Personal Information")
	default:
		return string(s)
	}
}

func (s Scope) DescriptionP(p *Printer) string {
	switch s {
	case ScopeBasic:
		return p.Sprintf("Read your Uid (login name) and Nickname")
	case ScopeEmail:
		return p.Sprintf("Read your Email address")
	case ScopeOpenID:
		return p.Sprintf("Read your ID Token after authenticated")
	case ScopeProfile:
		return p.Sprintf("Read your GivenName, Surname, BirthDate, etc.")
	default:
		return string(s)
	}
}

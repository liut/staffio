package web

// IDToken The ID Token represents a JWT passed to the client as part of the token response.
//
// https://openid.net/specs/openid-connect-core-1_0.html#IDToken
type IDToken struct {
	Issuer     string `json:"iss"` // REQUIRED. Issuer Identifier for the Issuer of the response.
	UserID     string `json:"sub"` // REQUIRED. Subject Identifier.
	ClientID   string `json:"aud"` // REQUIRED. Audience(s) that this ID Token is intended for.
	Expiration int64  `json:"exp"` // REQUIRED. Expiration time on or after which the ID Token
	IssuedAt   int64  `json:"iat"` // REQUIRED. Time at which the JWT was issued.

	Nonce string `json:"nonce,omitempty"` // Non-manditory fields MUST be "omitempty"

	// Custom claims supported by this server.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims

	Email         string `json:"email,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`

	UID        string `json:"uid,omitempty"`
	Name       string `json:"name,omitempty"`
	FamilyName string `json:"family_name,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	Locale     string `json:"locale,omitempty"`

	BirthDate   string `json:"birthdate,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Picture     string `json:"picture,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

func (z *IDToken) ToMap() map[string]any {
	if len(z.Issuer) == 0 || len(z.UserID) == 0 || len(z.ClientID) == 0 ||
		z.Expiration == 0 || z.IssuedAt == 0 {
		return nil
	}
	return map[string]any{
		"iss": z.Issuer,
		"sub": z.UserID,
		"aud": z.ClientID,
		"exp": z.Expiration,
		"iat": z.IssuedAt,

		"nonce": z.Nonce,
		"scope": "openid",
	}
}

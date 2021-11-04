package web

// IDToken The ID Token represents a JWT passed to the client as part of the token response.
//
// https://openid.net/specs/openid-connect-core-1_0.html#IDToken
type IDToken struct {
	Issuer     string `json:"iss"`
	UserID     string `json:"sub"`
	ClientID   string `json:"aud"`
	Expiration int64  `json:"exp"`
	IssuedAt   int64  `json:"iat"`

	Nonce string `json:"nonce,omitempty"` // Non-manditory fields MUST be "omitempty"

	// Custom claims supported by this server.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims

	Email         string `json:"email,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`

	Name       string `json:"name,omitempty"`
	FamilyName string `json:"family_name,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	Locale     string `json:"locale,omitempty"`

	BirthDate   string `json:"birthdate,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

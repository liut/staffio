package oidc

// IdToken defined from https://openid.net/specs/openid-connect-core-1_0.html#IDToken
// example:
//
//	{
//	 "iss": "https://server.example.com",
//	 "sub": "24400320",
//	 "aud": "s6BhdRkqt3",
//	 "nonce": "n-0S6_WzA2Mj",
//	 "exp": 1311281970,
//	 "iat": 1311280970,
//	 "auth_time": 1311280969,
//	 "acr": "urn:mace:incommon:iap:silver"
//	}
type IDToken struct {
	// REQUIRED. Issuer Identifier for the Issuer of the response.
	Iss string `json:"iss"`
	// REQUIRED. Subject Identifier. A locally unique and never reassigned identifier UID
	Sub string `json:"sub"`
	// REQUIRED. Audience(s) that this ID Token is intended for.
	Aud string `json:"aud"`
	// REQUIRED. Time at which the JWT was issued.
	Iat int64 `json:"iat"`
}

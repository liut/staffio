package oidc

import (
	"fmt"
	"strings"
)

type Discovery struct {
	Issuer                                 string   `json:"issuer"`
	AuthorizationEndpoint                  string   `json:"authorization_endpoint"`
	TokenEndpoint                          string   `json:"token_endpoint"`
	UserinfoEndpoint                       string   `json:"userinfo_endpoint"`
	JwksUri                                string   `json:"jwks_uri"`
	RegistrationEndpoint                   string   `json:"registration_endpoint,omitempty"`
	ResponseTypesSupported                 []string `json:"response_types_supported,omitempty"`
	ResponseModesSupported                 []string `json:"response_modes_supported,omitempty"`
	GrantTypesSupported                    []string `json:"grant_types_supported,omitempty"`
	SubjectTypesSupported                  []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported       []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                        []string `json:"scopes_supported"`
	ClaimsSupported                        []string `json:"claims_supported"`
	RequestParameterSupported              bool     `json:"request_parameter_supported,omitempty"`
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported,omitempty"`
	EndSessionEndpoint                     string   `json:"end_session_endpoint"`
}

func DiscoveryWith(uriPrefix string) Discovery {
	uriPrefix = strings.TrimRight(uriPrefix, "/")
	od := Discovery{
		Issuer:                                 uriPrefix,
		AuthorizationEndpoint:                  fmt.Sprintf("%s/authorize", uriPrefix),
		TokenEndpoint:                          fmt.Sprintf("%s/api/token", uriPrefix),
		UserinfoEndpoint:                       fmt.Sprintf("%s/api/info/generic", uriPrefix),
		JwksUri:                                fmt.Sprintf("%s/.well-known/jwks", uriPrefix), // TODO:
		ResponseTypesSupported:                 []string{"code", "token", "id_token", "code token", "code id_token", "token id_token", "code token id_token"},
		ResponseModesSupported:                 []string{"query"},
		GrantTypesSupported:                    []string{"password", "authorization_code"},
		SubjectTypesSupported:                  []string{"public"},
		IdTokenSigningAlgValuesSupported:       []string{"RS256"},
		ScopesSupported:                        []string{"openid", "email", "profile", "address", "phone"},
		ClaimsSupported:                        []string{"iss", "ver", "sub", "aud", "iat", "exp", "id", "type", "displayName", "avatar", "email", "phone"},
		RequestParameterSupported:              true,
		RequestObjectSigningAlgValuesSupported: []string{"HS256", "HS384", "HS512"},
		EndSessionEndpoint:                     fmt.Sprintf("%s/api/logout", uriPrefix),
	}
	return od
}

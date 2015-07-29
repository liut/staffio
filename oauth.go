package main

import (
	"github.com/RangelReale/osin"
)

func NewServerConfig() *osin.ServerConfig {
	return &osin.ServerConfig{
		AuthorizationExpiration: 900,
		AccessExpiration:        3600,
		TokenType:               "bearer",
		AllowedAuthorizeTypes: osin.AllowedAuthorizeType{
			osin.CODE,
			osin.TOKEN,
		},
		AllowedAccessTypes: osin.AllowedAccessType{
			osin.AUTHORIZATION_CODE,
			osin.IMPLICIT,
			osin.REFRESH_TOKEN,
			osin.PASSWORD,
			osin.CLIENT_CREDENTIALS,
		},
		ErrorStatusCode:           200,
		AllowClientSecretInParams: false,
		AllowGetAccessRequest:     false,
	}
}

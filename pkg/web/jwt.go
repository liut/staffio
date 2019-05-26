package web

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/openshift/osin"

	"github.com/liut/staffio/pkg/settings"
)

// JWT access token generator
type AccessTokenGenJWT struct {
	Key []byte
}

func (c *AccessTokenGenJWT) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	// generate JWT access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cid": data.Client.GetId(),
		"exp": data.ExpireAt().Unix(),
		"sub": data.UserData.(string),
	})

	accesstoken, err = token.SignedString(c.Key)
	if err != nil {
		return "", "", err
	}

	if !generaterefresh {
		return
	}

	// generate JWT refresh token
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cid": data.Client.GetId(),
	})

	refreshtoken, err = token.SignedString(c.Key)
	if err != nil {
		return "", "", err
	}
	return
}

func getTokenGenJWT() (tokenGen osin.AccessTokenGen, err error) {
	var (
		hmacKey []byte
	)

	hmacKey, err = jwt.DecodeSegment(settings.Current.TokenGenKey)
	if err != nil {
		log.Printf("ERROR: key %s\n", err)
		return
	}

	tokenGen = &AccessTokenGenJWT{Key: hmacKey}

	return
}

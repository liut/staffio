package web

import (
	"fmt"

	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"

	. "lcgc/platform/staffio/settings"
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

	hmacKey, err = jwt.DecodeSegment(Settings.TokenGen.Key)
	if err != nil {
		fmt.Printf("ERROR: key %s\n", err)
		return
	}
	fmt.Printf("%v\n", hmacKey)

	tokenGen = &AccessTokenGenJWT{Key: hmacKey}

	return
}

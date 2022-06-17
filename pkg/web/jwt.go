package web

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openshift/osin"

	"github.com/liut/staffio/pkg/models/oauth"
	"github.com/liut/staffio/pkg/settings"
)

// AccessTokenGenJWT JWT access token generator
type AccessTokenGenJWT struct {
	Key []byte
}

func (c *AccessTokenGenJWT) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	// generate JWT access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cid": data.Client.GetId(),
		"exp": data.ExpireAt().Unix(),
		"sub": oauth.StringFromMeta(data.UserData, "uid"),
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
		logger().Warnw("getTokenGenJWT fail", "err", err)
		return
	}

	tokenGen = &AccessTokenGenJWT{Key: hmacKey}

	return
}

// LoadPrivateKey loads a private key from PEM/DER data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	return nil, fmt.Errorf("parse error, got '%s', '%s', '%s'", err0, err1, err2)
}

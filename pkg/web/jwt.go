package web

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/go-osin/osin"
	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/liut/staffio/pkg/models/oauth"
	"github.com/liut/staffio/pkg/settings"
)

type MapGetter interface {
	ToMap() map[string]any
}

type TokenGenerator interface {
	osin.AccessTokenGen
	GenerateIDToken(mg MapGetter) (string, error)
	getJSONWebKey() jose.JSONWebKey
}

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

func (c *AccessTokenGenJWT) GenerateIDToken(mg MapGetter) (string, error) {
	mc := mg.ToMap()
	if mc == nil {
		return "", fmt.Errorf("invalid map")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(mc))

	idtoken, err := token.SignedString(c.Key)
	if err != nil {
		return "", err
	}
	return idtoken, nil
}

func getTokenGenJWT() (tokenGen TokenGenerator, err error) {
	var (
		hmacKey []byte
	)

	hmacKey, err = base64.RawURLEncoding.DecodeString(settings.Current.TokenGenKey)
	if err != nil {
		logger().Warnw("getTokenGenJWT fail", "err", err)
		return
	}

	tokenGen = &AccessTokenGenJWT{Key: hmacKey}

	return
}

func (c *AccessTokenGenJWT) getJSONWebKey() jose.JSONWebKey {

	var jwk jose.JSONWebKey
	jwk.Key = c.Key
	jwk.Certificates = []*x509.Certificate{}
	jwk.KeyID = "cust"
	jwk.Algorithm = jwt.SigningMethodHS256.Name
	jwk.Use = "sig"

	return jwk
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

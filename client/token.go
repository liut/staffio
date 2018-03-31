package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/oauth2"
)

type InfoToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	Me           User      `json:"me,omitempty"`
}

func (tok *InfoToken) GetExpiry() time.Time {
	return time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
}

type infoError struct {
	Code    string `json:"error,omitempty"`
	Message string `json:"error_description,omitempty"`
}

func (e *infoError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func requestInfoToken(tok *oauth2.Token) (*InfoToken, error) {
	client := conf.Client(oauth2.NoContext, tok)
	info, err := client.Get(infoUrl)
	if err != nil {
		return nil, err
	}
	defer info.Body.Close()
	data, err := ioutil.ReadAll(info.Body)
	if err != nil {
		log.Printf("read err %s", err)
		return nil, err
	}
	log.Print(string(data))

	infoErr := &infoError{}
	if e := json.Unmarshal(data, infoErr); e != nil {
		return nil, e
	}

	if infoErr.Code != "" {
		return nil, infoErr
	}

	var token = &InfoToken{}
	err = json.Unmarshal(data, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

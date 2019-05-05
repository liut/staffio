package oauth

import (
	"time"
)

type Client struct {
	Id                   uint        `json:"_id,omitempty"`
	Name                 string      `json:"name"`
	Code                 string      `json:"code,omitempty"`
	Secret               string      `json:"-"`
	RedirectUri          string      `json:"uri"`
	UserData             interface{} `json:"-"`
	CreatedAt            time.Time   `json:"created,omitempty"`
	AllowedGrantTypes    []string    `json:"grant_types,omitempty"`
	AllowedResponseTypes []string    `json:"response_types,omitempty"`
	AllowedScopes        []string    `json:"scopes,omitempty"`
}

func (c *Client) GetId() string {
	return c.Code
}

func (c *Client) GetSecret() string {
	return c.Secret
}

func (c *Client) GetRedirectUri() string {
	return c.RedirectUri
}

func (c *Client) GetUserData() interface{} {
	return c.UserData
}

func NewClient(name, code, secret, redirectUri string) *Client {
	return &Client{
		Name:              name,
		Code:              code,
		Secret:            secret,
		RedirectUri:       redirectUri,
		CreatedAt:         time.Now(),
		AllowedGrantTypes: []string{"authorization_code", "refresh_token"},
		AllowedScopes:     []string{"basic"},
	}
}

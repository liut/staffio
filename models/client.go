package models

import (
	"time"
)

type Client struct {
	Id          uint `json:"_id,omitempty"`
	Name        string
	Code        string `json:"code,omitempty"`
	Secret      string
	RedirectUri string
	UserData    interface{}
	Created     time.Time `json:"created,omitempty"`
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

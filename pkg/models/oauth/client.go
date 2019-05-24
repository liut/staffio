package oauth

import (
	"time"

	"github.com/liut/staffio/pkg/models/types"
)

type StringSlice = types.StringSlice

// Client of oauth2 app
type Client struct {
	ID                   uint        `json:"id,omitempty"`
	Name                 string      `json:"name"`
	Code                 string      `json:"code,omitempty"`
	Secret               string      `json:"secret,omitempty"`
	RedirectURI          string      `json:"redirect_uri" db:"redirect_uri"`
	UserData             interface{} `json:"-" db:"userdata"`
	CreatedAt            time.Time   `json:"created,omitempty" db:"created"`
	AllowedGrantTypes    StringSlice `json:"grant_types,omitempty" db:"grant_types" `
	AllowedResponseTypes StringSlice `json:"response_types,omitempty" db:"response_types"`
	AllowedScopes        StringSlice `json:"scopes,omitempty" db:"scopes"`
}

// GetId osin.Client.GetId
func (c *Client) GetId() string {
	return c.Code
}

// GetSecret osin.Client.GetSecret
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetRedirectUri osin.Client.GetRedirectUri
func (c *Client) GetRedirectUri() string {
	return c.RedirectURI
}

// GetUserData osin.Client.GetUserData
func (c *Client) GetUserData() interface{} {
	return c.UserData
}

// NewClient build a client
func NewClient(name, code, secret, redirectURI string) *Client {
	return &Client{
		Name:              name,
		Code:              code,
		Secret:            secret,
		RedirectURI:       redirectURI,
		CreatedAt:         time.Now(),
		AllowedGrantTypes: []string{"authorization_code", "refresh_token"},
		AllowedScopes:     []string{"basic"},
	}
}

// ClientSpec 查询参数
type ClientSpec struct {
	Page   int      `json:"page,omitempty" form:"page"`
	Limit  int      `json:"limit,omitempty" form:"limit"`
	Orders []string `json:"order,omitempty" form:"order"`
	Total  int      `json:"total,omitempty"` // for set value

	CountOnly bool `json:"count,omitempty" form:"count"`
}

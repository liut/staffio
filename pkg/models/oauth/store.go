package oauth

import (
	"github.com/openshift/osin"
)

type OSINStore interface {
	osin.Storage

	LoadClients(spec *ClientSpec) ([]Client, error)
	CountClients() uint
	GetClientWithCode(code string) (*Client, error)
	GetClientWithID(id int) (*Client, error)
	SaveClient(client *Client) error
	RemoveClient(code string) error

	LoadScopes() (scopes []Scope, err error)
	IsAuthorized(clientID, username string) bool
	SaveAuthorized(clientID, username string) error
}

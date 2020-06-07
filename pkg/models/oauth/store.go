package oauth

import (
	"github.com/openshift/osin"
)

// OSINStore ...
type OSINStore interface {
	osin.Storage

	LoadClient(id string) (*Client, error)
	LoadClients(spec *ClientSpec) ([]Client, error)
	CountClients() uint
	SaveClient(client *Client) error
	RemoveClient(id string) error

	LoadScopes() (scopes []Scope, err error)
	IsAuthorized(clientID, username string) bool
	SaveAuthorized(clientID, username string) error
}

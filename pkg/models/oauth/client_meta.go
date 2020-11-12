package oauth

import (
	"database/sql/driver"
	"encoding/json"
)

var (
	defaultGrantTypes    = []string{"authorization_code", "password", "refresh_token"}
	defaultResponseTypes = []string{"code", "token"}
	defaultScopes        = []string{"basic"}
	defaultClientMeta    = ClientMeta{
		Name:          "",
		GrantTypes:    defaultGrantTypes,
		ResponseTypes: defaultResponseTypes,
		Scopes:        defaultScopes,
	}
)

// ClientMeta ...
type ClientMeta struct {
	Name          string   `json:"name,omitempty"`
	GrantTypes    []string `json:"grant_types,omitempty"`    // AllowedGrantTypes
	ResponseTypes []string `json:"response_types,omitempty"` // AllowedResponseTypes
	Scopes        []string `json:"scopes,omitempty"`         // AllowedScopes
}

// Scan implements the sql.Scanner interface.
func (m *ClientMeta) Scan(value interface{}) (err error) {
	switch data := value.(type) {
	case ClientMeta:
		*m = data
	case []byte:
		err = json.Unmarshal(data, m)
	case string:
		err = json.Unmarshal([]byte(data), m)
	}
	return
}

// Value implements the driver.Valuer interface.
func (m ClientMeta) Value() (driver.Value, error) {
	return json.Marshal(m)
}

package oauth

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// vars
var (
	ErrInvalidJSON = errors.New("Invalid JSON")
)

// JSONKV ...
type JSONKV map[string]interface{}

// ToJSONKV ...
func ToJSONKV(src interface{}) (JSONKV, error) {
	switch s := src.(type) {
	case JSONKV:
		return s, nil
	case map[string]interface{}:
		return JSONKV(s), nil
	}
	return nil, ErrInvalidJSON
}

// WithKey ...
func (m JSONKV) WithKey(key string) (v interface{}) {
	var ok bool
	if v, ok = m[key]; ok {
		return
	}
	return
}

// Scan implements the sql.Scanner interface.
func (m *JSONKV) Scan(value interface{}) (err error) {
	switch data := value.(type) {
	case JSONKV:
		*m = data
	case map[string]interface{}:
		*m = JSONKV(data)
	case []byte:
		err = json.Unmarshal(data, m)
	case string:
		err = json.Unmarshal([]byte(data), m)
	}
	return
}

// Value implements the driver.Valuer interface.
func (m JSONKV) Value() (driver.Value, error) {
	return json.Marshal(m)
}

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

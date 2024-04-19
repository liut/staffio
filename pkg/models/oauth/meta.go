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
	case string:
		return JSONKV{"_": s}, nil
	}
	return JSONKV{}, ErrInvalidJSON
}

// WithKey ...
func (m JSONKV) WithKey(key string) (v interface{}) {
	var ok bool
	if v, ok = m[key]; ok { //nolint
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

// StringFromMeta ...
func StringFromMeta(kv interface{}, key string) string {
	if m, err := ToJSONKV(kv); err == nil {
		if v, ok := m[key]; ok {
			return v.(string)
		}
	}
	return ""
}

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
type JSONKV map[string]any

// ToJSONKV ...
func ToJSONKV(src any) (JSONKV, error) {
	switch s := src.(type) {
	case JSONKV:
		return s, nil
	case map[string]any:
		return JSONKV(s), nil
	case string:
		return JSONKV{"_": s}, nil
	}
	return JSONKV{}, ErrInvalidJSON
}

// Get ...
func (m JSONKV) Get(key string) (v any, ok bool) {
	v, ok = m[key]
	return
}

func (m JSONKV) GetInt(key string) int {
	if v, ok := m[key]; ok {
		switch z := v.(type) {
		case float64:
			return int(z)
		case int:
			return z
		}
	}
	return 0
}

func (m JSONKV) GetStr(key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Scan implements the sql.Scanner interface.
func (m *JSONKV) Scan(value any) (err error) {
	switch data := value.(type) {
	case JSONKV:
		*m = data
	case map[string]any:
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
func StringFromMeta(kv any, key string) string {
	if m, err := ToJSONKV(kv); err == nil {
		if v, ok := m[key]; ok {
			return v.(string)
		}
	}
	return ""
}

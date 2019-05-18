package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var ErrAssertion = errors.New("type assertion to []byte failed")

type StringSlice []string

// Value Make the StringSlice implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the StringSlice.
func (a StringSlice) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Make the StringSlice implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the StringSlice.
func (a *StringSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return ErrAssertion
	}

	return json.Unmarshal(b, &a)
}

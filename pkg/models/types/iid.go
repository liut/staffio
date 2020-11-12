package types

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

// IID Integer ID
type IID uint64

// Bytes ...
func (z IID) Bytes() []byte {
	var bInt big.Int
	return bInt.SetUint64(uint64(z)).Bytes()
}

// String ...
func (z IID) String() string {
	var bInt big.Int
	return bInt.SetUint64(uint64(z)).Text(36)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (z IID) MarshalText() ([]byte, error) {
	b := []byte(z.String())
	return b, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (z *IID) UnmarshalText(data []byte) (err error) {
	var id IID
	id, err = ParseID(string(data))
	*z = id
	return
}

// ParseID ...
func ParseID(s string) (IID, error) {
	var id uint64
	var bI big.Int
	if i, ok := bI.SetString(s, 36); ok {
		id = i.Uint64()
	} else {
		return 0, fmt.Errorf("invalid id %q", s)
	}
	return IID(id), nil
}

// Scan implements of database/sql.Scanner
func (z *IID) Scan(src interface{}) (err error) {
	switch s := src.(type) {
	case string:
		return z.UnmarshalText([]byte(s))
	case []byte:
		return z.UnmarshalText(s)
	}
	return fmt.Errorf("'%v' is invalid IID", src)
}

// Value implements of database/sql/driver.Valuer
func (z IID) Value() (driver.Value, error) {
	return z.String(), nil
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

package auth

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *User) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "u"
	o = append(o, 0x84, 0xa1, 0x75)
	o = msgp.AppendString(o, z.UID)
	// string "n"
	o = append(o, 0xa1, 0x6e)
	o = msgp.AppendString(o, z.Name)
	// string "p"
	o = append(o, 0xa1, 0x70)
	o = msgp.AppendString(o, z.Privileges)
	// string "h"
	o = append(o, 0xa1, 0x68)
	o = msgp.AppendInt64(o, z.LastHit)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *User) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "u":
			z.UID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "n":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "p":
			z.Privileges, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "h":
			z.LastHit, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *User) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.UID) + 2 + msgp.StringPrefixSize + len(z.Name) + 2 + msgp.StringPrefixSize + len(z.Privileges) + 2 + msgp.Int64Size
	return
}

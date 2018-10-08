package client

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *User) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "u":
			z.Uid, err = dc.ReadString()
			if err != nil {
				return
			}
		case "n":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "h":
			z.LastHit, err = dc.ReadInt64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z User) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "u"
	err = en.Append(0x83, 0xa1, 0x75)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Uid)
	if err != nil {
		return
	}
	// write "n"
	err = en.Append(0xa1, 0x6e)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "h"
	err = en.Append(0xa1, 0x68)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.LastHit)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z User) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "u"
	o = append(o, 0x83, 0xa1, 0x75)
	o = msgp.AppendString(o, z.Uid)
	// string "n"
	o = append(o, 0xa1, 0x6e)
	o = msgp.AppendString(o, z.Name)
	// string "h"
	o = append(o, 0xa1, 0x68)
	o = msgp.AppendInt64(o, z.LastHit)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *User) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "u":
			z.Uid, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "n":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
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
func (z User) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.Uid) + 2 + msgp.StringPrefixSize + len(z.Name) + 2 + msgp.Int64Size
	return
}

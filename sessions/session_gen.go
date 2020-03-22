package sessions

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (s *Session) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "v":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if s.Values == nil && msz > 0 {
				s.Values = make(map[string]interface{}, msz)
			} else if len(s.Values) > 0 {
				for key := range s.Values {
					delete(s.Values, key)
				}
			}
			for msz > 0 {
				msz--
				var xvk string
				var bzg interface{}
				xvk, err = dc.ReadString()
				if err != nil {
					return
				}
				bzg, err = dc.ReadIntf()
				if err != nil {
					return
				}
				s.Values[xvk] = bzg
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
func (s *Session) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "v"
	err = en.Append(0x81, 0xa1, 0x76)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(s.Values)))
	if err != nil {
		return
	}
	for xvk, bzg := range s.Values {
		err = en.WriteString(xvk)
		if err != nil {
			return
		}
		err = en.WriteIntf(bzg)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (s *Session) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, s.Msgsize())
	// map header, size 1
	// string "v"
	o = append(o, 0x81, 0xa1, 0x76)
	o = msgp.AppendMapHeader(o, uint32(len(s.Values)))
	for xvk, bzg := range s.Values {
		o = msgp.AppendString(o, xvk)
		o, err = msgp.AppendIntf(o, bzg)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (s *Session) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "v":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if s.Values == nil && msz > 0 {
				s.Values = make(map[string]interface{}, msz)
			} else if len(s.Values) > 0 {
				for key := range s.Values {
					delete(s.Values, key)
				}
			}
			for msz > 0 {
				var xvk string
				var bzg interface{}
				msz--
				xvk, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bzg, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				s.Values[xvk] = bzg
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

func (s *Session) Msgsize() (sz int) {
	sz = 1 + 2 + msgp.MapHeaderSize
	if s.Values != nil {
		for xvk, bzg := range s.Values {
			_ = bzg
			sz += msgp.StringPrefixSize + len(xvk) + msgp.GuessSize(bzg)
		}
	}
	return
}

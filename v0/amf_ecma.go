package v0

import (
	"fmt"
	"bytes"
	"encoding/binary"
)

type AmfECMA map[string]interface{}

func NewAmfECMA() *AmfECMA {
	ecma := make(AmfECMA)
	return &ecma
}

func(c *AmfCodec) AmfECMAEncode(val *AmfECMA) ([]byte, error) {
	if val == nil {
		return nil, fmt.Errorf("Empty AmfECMA")
	}
	if id, err := c.GetId(val); err == nil {
		return AmfRefEncode(id), nil
	}

	res := []byte{AMF_ECMA}
	datalen := len(*val)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(datalen))
	res = append(res, buf...)
	for key, value := range *val {
		keybytes := AmfUTF8Encode(key)
		valuebytes, err := c.Encode(value)
		if err != nil {
			return nil, err
		}
		res = append(res, keybytes...)
		res = append(res, valuebytes...)
	}
	res = append(res, AmfObjEnd()...)
	c.Set(val)
	return res, nil
}

func(c *AmfCodec) AmfECMADecode(data []byte) (*AmfECMA, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("Empty data")
	}
	if data[0] == AMF_REFERENCE {
		id, _, err := AmfRefDecode(data)
		if err != nil {
			return nil, 0, err
		}
		if obj, err := c.GetObj(id); err == nil {
			if ecma, ok := obj.(*AmfECMA); ok {
				return ecma, 3, nil
			} else {
				return nil, 0, fmt.Errorf("Value in table is not an AmfECMA")
			}
		} else {
			return nil, 0, err
		}
	}
	if data[0] != AMF_ECMA {
		return nil, 0, fmt.Errorf("Invalid type")
	}
	data = data[1:]
	totalConsumed := 1
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfECMA")
	}
	// ECMA's length is more often than not, not an accurate representation of the data's length
	_ = binary.BigEndian.Uint32(data[:4])
	data = data[4:]
	totalConsumed += 4
	end := AmfObjEnd()

	obj := NewAmfECMA()
	for {
		if len(data) < 3 {
			return nil, 0, fmt.Errorf("No OBJ_END marker found. ECMA may be corrupted")
		}
		if bytes.Equal(data[:3], end) {
			c.Set(obj)
			return obj, totalConsumed + 3, nil
		}
		key, cnt, err := AmfUTF8Decode(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt

		value, cnt, err := c.Decode(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt

		(*obj)[key] = value
	}
}

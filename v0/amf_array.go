package v0

import (
	"fmt"
	"encoding/binary"
)

type AmfArray []interface{}

func NewAmfArray() *AmfArray {
	return &AmfArray{}
}

func(arr *AmfArray) Add(value interface{}) {
	*arr = append(*arr, value)
}

func(c *AmfCodec) AmfArrayEncode(arr *AmfArray) ([]byte, error) {
	if arr == nil {
		return nil, fmt.Errorf("Empty AmfArray")
	}
	// ref
	if id, err := c.GetId(arr); err == nil {
		return AmfRefEncode(id), nil
	}
	
	res := []byte{AMF_STRICTARR}
	size := uint32(len(*arr))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, size)
	res = append(res, buf...)
	for _, value := range *arr {
		valuebytes, err := c.Encode(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valuebytes...)
	}
	c.Set(arr)
	return res, nil
}

func(c *AmfCodec) AmfArrayDecode(data []byte) (*AmfArray, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("Empty data")
	}
	if data[0] == AMF_REFERENCE {
		id, _, err := AmfRefDecode(data)
		if err != nil {
			return nil, 0, err
		}
		if obj, err := c.GetObj(id); err == nil {
			if arr, ok := obj.(*AmfArray); ok {
				return arr, 3, nil
			} else {
				return nil, 0, fmt.Errorf("Value in table is not an AmfArray")
			}
		} else {
			return nil, 0, err
		}
	}
	if data[0] != AMF_STRICTARR {
		return nil, 0, fmt.Errorf("Invalid type")
	}
	data = data[1:]
	totalConsumed := 1
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfArray")
	}
	datalen := binary.BigEndian.Uint32(data[:4])
	data = data[4:]
	totalConsumed += 4

	arr := NewAmfArray()
	for i := uint32(0); i < datalen; i++ {
		value, cnt, err := c.Decode(data)
		if err != nil {
			return nil, 0, err
		}
		*arr = append(*arr, value)
		data = data[cnt:]
		totalConsumed += cnt
	}
	c.Set(arr)
	return arr, totalConsumed, nil
}

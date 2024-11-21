package v0

import (
	"fmt"
	"bytes"
)

type AmfObj struct {
	Name string
	Props []AmfObjProp
}

type AmfObjProp struct {
	Key string
	Value interface{}
}

func(obj *AmfObj) AddProp(key string, val interface{}) {
	obj.Props = append(obj.Props, AmfObjProp{key, val})
}

func NewAmfObj() *AmfObj {
	return &AmfObj{
		Name: "",
		Props: make([]AmfObjProp, 0),
	}
}

func AmfObjEnd() []byte {
	return []byte{0x00, 0x00, 0x09}
}

// Used to encode an object. Because object passed in is strongly typed, we encode into 
// the TypedObject format.
// Because we're using `reflect`, all fields must be exported, the unexported fields will be ignored.
func(c *AmfCodec) AmfObjEncode(obj *AmfObj) ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Empty AmfObj")
	}

	// Found in ref table, encode as reference
	if id, err := c.GetId(obj); err == nil {
		return AmfRefEncode(id), nil
	}
	
	res := []byte{AMF_OBJECT}
	if obj.Name != "" {
		namebytes := AmfUTF8Encode(obj.Name)
		res[0] = AMF_TYPEDOBJ
		res = append(res, namebytes...)
	}
	for _, prop := range obj.Props {
		key := AmfUTF8Encode(prop.Key)
		value, err := c.Encode(prop.Value)
		if err != nil {
			return nil, err
		}
		res = append(res, key...)
		res = append(res, value...)
	}
	res = append(res, AmfObjEnd()...)
	c.Set(obj)
	return res, nil
}

func(c *AmfCodec) AmfObjDecode(data []byte) (*AmfObj, int, error) {
	if len(data) < 1{
		return nil, 0, fmt.Errorf("Not enough data to decode AmfObj")
	}
	// reference
	if data[0] == AMF_REFERENCE {
		id, cnt, err := AmfRefDecode(data)
		if err != nil {
			return nil, 0, err
		}
		res, err := c.GetObj(id)
		if err != nil {
			return nil, 0, err
		}
		if res, ok := res.(*AmfObj); ok {
			return res, cnt, nil
		} else {
			return nil, 0, fmt.Errorf("Value in table is not an AmfObj")
		}
	}

	obj := NewAmfObj()
	totalConsumed := 0
	// raw
	if data[0] == AMF_OBJECT {
		data = data[1:]
		totalConsumed = 1
	} else if data[0] == AMF_TYPEDOBJ {
		data = data[1:]
		totalConsumed = 1
		name, cnt, err := AmfUTF8Decode(data)
		if err != nil {
			return nil, 0, err
		}
		obj.Name = name
		data = data[cnt:]
		totalConsumed += cnt
	} else {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_OBJECT, data[0])
	}
	objend := AmfObjEnd()

	for {
		datalen := len(data)
		if datalen < 3 {
			return nil, 0, fmt.Errorf("Not enough data to decode AmfObj")
		}
		if bytes.Equal(data[:3], objend) {
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

		obj.Props = append(obj.Props, AmfObjProp{key, value})
	}
}

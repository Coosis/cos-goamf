package v3

import (
	"fmt"
)

type AmfVectorObj struct {
	TypeName string
	FixedLen bool
	Data []interface{}
}

func EmptyAmfVectorObj() *AmfVectorObj {
	return &AmfVectorObj{
		TypeName: "*",
		FixedLen: false,
		Data: make([]interface{}, 0),
	}
}

func(codec *AmfCodec) AmfVectorObjEncode(
	obj *AmfVectorObj,
) ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Empty AmfVectorObj")
	}
	data := obj.Data
	if data == nil {
		return nil, fmt.Errorf("Empty AmfVectorObj data")
	}
	name := obj.TypeName
	fixedLen := obj.FixedLen

	if id, ok := codec.GetId(obj, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_VECTOR_OBJECT}, objref...), nil
	}

	length := uint32(len(data))
	U29o, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := U29o

	if fixedLen {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}
	typeName, err := codec.AmfStringEncodePayload(name)
	if err != nil {
		return nil, err
	}
	body = append(body, typeName...)
	for i := range length {
		obj, err := codec.AmfEncode(data[i])
		if err != nil {
			return nil, err
		}
		body = append(body, obj...)
	}
	codec.Append(obj, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_OBJECT}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorObjDecode(data []byte) (*AmfVectorObj, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorObject")
	}
	if data[0] != AMF_VECTOR_OBJECT {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_VECTOR_OBJECT, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29o, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29o & 1) == 0 {
		// o-ref
		id := U29o >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.(*AmfVectorObj); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a vector object")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		// u29o-value *(object-type-name object)
		length := U29o >> 1
		if uint32(len(data)) < length {
			return nil, 0, fmt.Errorf("Vector object is shorter than expected: %v", data)
		}
		obj := EmptyAmfVectorObj()
		obj.FixedLen = data[0] == 0x01
		data = data[1:]
		totalConsumed++

		Typename, cnt, err := codec.AmfStringDecodePayload(data)
		obj.TypeName = Typename
		if err != nil {
			return nil, 0, err
		}
		totalConsumed += cnt
		data = data[cnt:]
		for i := 0; i < int(length); i++ {
			val, cnt, err := codec.AmfDecode(data)
			if err != nil {
				return nil, 0, err
			}
			totalConsumed += cnt
			data = data[cnt:]
			obj.Data = append(obj.Data, val)
		}

		codec.Append(obj, COMPLEX_TABLE)
		return obj, totalConsumed, nil
	}
}

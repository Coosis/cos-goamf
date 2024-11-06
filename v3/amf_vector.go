package v3

import (
	"fmt"
	"math"
	"encoding/binary"
)

func(codec *AmfCodec) AmfVectorDoubleEncode(data *[]float64, fixedLen bool) ([]byte, error) {
	if id, ok := codec.GetId(data, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_VECTOR_DOUBLE}, objref...), nil
	}

	length := uint32(len(*data))
	U29d, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := U29d
	if fixedLen {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}

	for _, v := range *data {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(v))
		body = append(body, buf...)
	}

	codec.Append(data, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_DOUBLE}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorDoubleDecode(data []byte) (*[]float64, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorDouble")
	}
	if data[0] != AMF_VECTOR_DOUBLE {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_VECTOR_DOUBLE, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29d, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29d & 1) == 0 {
		// d-ref
		id := U29d >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.(*[]float64); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a vector double")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		// discard the "fixed-vector" flag
		if len(data) < 1 {
			return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorDouble")
		}
		data = data[1:]
		totalConsumed++

		// u29d-value *(double-value)
		length := U29d >> 1
		if uint32(len(data)) < length * 8 {
			return nil, 0, fmt.Errorf("Vector double is shorter than expected: %v", data)
		}
		res := make([]float64, length)
		for i := 0; i < int(length); i++ {
			buf := data[i*8:i*8+8]
			res[i] = math.Float64frombits(binary.BigEndian.Uint64(buf))
		}

		codec.Append(&res, COMPLEX_TABLE)
		return &res, totalConsumed, nil
	}
}

func(codec *AmfCodec) AmfVectorIntEncode(data *[]int32, fixedLen bool) ([]byte, error) {
	if id, ok := codec.GetId(data, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_VECTOR_INT}, objref...), nil
	}

	length := uint32(len(*data))
	U29i, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := U29i

	if fixedLen {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}
	for _, v := range *data {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		body = append(body, buf...)
	}

	codec.Append(data, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_INT}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorIntDecode(data []byte) (*[]int32, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorInt")
	}
	if data[0] != AMF_VECTOR_INT {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_VECTOR_INT, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29i, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29i & 1) == 0 {
		// i-ref
		id := U29i >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.(*[]int32); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a vector int")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		// discard the "fixed-vector" flag
		if len(data) < 1 {
			return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorInt")
		}
		data = data[1:]
		totalConsumed++

		// u29i-value *(int-value)
		length := U29i >> 1
		if uint32(len(data)) < length * 4 {
			return nil, 0, fmt.Errorf("Vector int is shorter than expected: %v", data)
		}
		res := make([]int32, length)
		for i := 0; i < int(length); i++ {
			buf := data[i*4:i*4+4]
			res[i] = int32(binary.BigEndian.Uint32(buf))
		}

		codec.Append(&res, COMPLEX_TABLE)
		return &res, totalConsumed, nil
	}
}

func(codec *AmfCodec) AmfVectorUintEncode(data *[]uint32, fixedLen bool) ([]byte, error) {
	if id, ok := codec.GetId(data, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_VECTOR_UINT}, objref...), nil
	}

	length := uint32(len(*data))
	U29i, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := U29i

	if fixedLen {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}
	for _, v := range *data {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		body = append(body, buf...)
	}

	codec.Append(data, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_UINT}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorUintDecode(data []byte) (*[]uint32, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorUint")
	}
	if data[0] != AMF_VECTOR_UINT {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_VECTOR_UINT, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29i, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29i & 1) == 0 {
		// i-ref
		id := U29i >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.(*[]uint32); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a vector uint")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		// discard the "fixed-vector" flag
		if len(data) < 1 {
			return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorUint")
		}
		data = data[1:]
		totalConsumed++

		// u29i-value *(uint-value)
		length := U29i >> 1
		if uint32(len(data)) < length * 4 {
			return nil, 0, fmt.Errorf("Vector uint is shorter than expected: %v", data)
		}
		res := make([]uint32, length)
		for i := 0; i < int(length); i++ {
			buf := data[i*4:i*4+4]
			res[i] = binary.BigEndian.Uint32(buf)
		}

		codec.Append(&res, COMPLEX_TABLE)
		return &res, totalConsumed, nil
	}
}

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
	marker AmfMarker,
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
		obj, err := codec.AmfEncode(data[i], marker)
		if err != nil {
			return nil, err
		}
		body = append(body, obj...)
	}
	codec.Append(obj, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_OBJECT}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorObjDecode(data []byte, marker AmfMarker) (*AmfVectorObj, int, error) {
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
			val, _, cnt, err := codec.AmfDecode(data)
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

package v3

import (
	"fmt"
)

// array-marker
// (
// 	U29O-ref |
// 	(
// 		U29A-value (
// 			UTF-8-empty | *(assoc-value)
// 			UTF-8-empty
// 		)
// 		*(value-type)
// 	)
// )

// TODO: Add type for dense elements and associative elements
type AmfArray struct {
	Dense []interface{}
	DenseMarker []AmfMarker

	Assoc map[string]interface{}
	AssocMarker map[string]AmfMarker
	AssocKeys []string
}

func EmptyAmfArray() *AmfArray {
	return &AmfArray{
		Dense: make([]interface{}, 0),
		DenseMarker: make([]AmfMarker, 0),
		Assoc: make(map[string]interface{}, 0),
		AssocMarker: make(map[string]AmfMarker, 0),
		AssocKeys: make([]string, 0),
	}
}

func(arr *AmfArray) AddDense(value interface{}, marker AmfMarker) {
	arr.Dense = append(arr.Dense, value)
	arr.DenseMarker = append(arr.DenseMarker, marker)
}

func(arr *AmfArray) AddAssoc(key string, value interface{}, marker AmfMarker) {
	arr.Assoc[key] = value
	arr.AssocMarker[key] = marker
	arr.AssocKeys = append(arr.AssocKeys, key)
}

func(codec *AmfCodec) AmfArrayEncode(arr *AmfArray) ([]byte, error) {
	if arr == nil {
		return nil, fmt.Errorf("Empty AmfArray")
	}
	if id, ok := codec.GetId(arr, COMPLEX_TABLE); ok {
		num, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_ARRAY}, num...), nil
	}

	// U29A-value
	size := uint32(len(arr.Dense))
	num, err := AmfIntEncodePayload(size << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := num
	for _, key := range arr.AssocKeys {
		value := arr.Assoc[key]
		keyEncoded, err := codec.AmfStringEncodePayload(key)
		if err != nil {
			return nil, err
		}
		valueEncoded, err := codec.AmfEncode(value, arr.AssocMarker[key])
		if err != nil {
			return nil, err
		}
		body = append(body, keyEncoded...)
		body = append(body, valueEncoded...)
	}
	body = append(body, 0x01)
	for i, value := range arr.Dense {
		encoded, err := codec.AmfEncode(value, arr.DenseMarker[i])
		if err != nil {
			return nil, err
		}
		body = append(body, encoded...)
	}
	codec.Append(arr, COMPLEX_TABLE)
	return append([]byte{AMF_ARRAY}, body...), nil
}

func(codec *AmfCodec) AmfArrayDecode(data []byte) (*AmfArray, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("Empty byte array")
	}
	if data[0] != AMF_ARRAY {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_ARRAY, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	num, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	totalConsumed += cnt
	data = data[cnt:]

	// ref
	if (num & 0x01) == 0 {
		id := num >> 1
		if obj, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if obj, ok := obj.(*AmfArray); ok {
				return obj, cnt, nil
			}
		}
		return nil, 0, fmt.Errorf("Invalid ref id: %v", id)
	}

	// U29A-value
	arr := EmptyAmfArray()
	size := num >> 1
	for {
		// assoc-value
		key, cnt, err := codec.AmfStringDecodePayload(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt
		if key == "" {
			break
		}
		value, marker, cnt, err := codec.AmfDecode(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt
		arr.AddAssoc(key, value, marker)
	}

	// value-type
	for i := 0; i < int(size); i++ {
		value, marker, cnt, err := codec.AmfDecode(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt
		arr.AddDense(value, marker)
	}

	codec.Append(arr, COMPLEX_TABLE)
	return arr, totalConsumed, nil
}

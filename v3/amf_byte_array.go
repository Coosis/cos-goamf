package v3

import (
	"fmt"
)

func(codec *AmfCodec) AmfByteArrayEncode(value *[]byte) ([]byte, error) {
	if id, ok := codec.GetId(value, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_BYTE_ARRAY}, objref...), nil
	}

	length := uint32(len(*value))
	U29b, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	codec.Append(value, COMPLEX_TABLE)
	body := append(U29b, *value...)
	res := append([]byte{AMF_BYTE_ARRAY}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfByteArrayDecode(data []byte) ([]byte, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfByteArray")
	}
	if data[0] != AMF_BYTE_ARRAY {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_BYTE_ARRAY, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29b, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29b & 1) == 0 {
		// o-ref
		id := U29b >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.([]byte); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a byte array")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		// u29b-value *(byte)
		length := U29b >> 1
		if uint32(len(data)) < length {
			return nil, 0, fmt.Errorf("Byte array is shorter than expected: %v", data)
		}
		return data[:length], totalConsumed, nil
	}
}

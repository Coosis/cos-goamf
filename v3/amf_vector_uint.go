package v3

import (
	"fmt"
	"encoding/binary"
)

type AmfVectorUint struct {
	FixedLen bool
	Data []uint32
}

func(codec *AmfCodec) AmfVectorUintEncode(vec *AmfVectorUint) ([]byte, error) {
	if id, ok := codec.GetId(vec, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_VECTOR_UINT}, objref...), nil
	}

	length := uint32(len(vec.Data))
	U29i, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := U29i

	if vec.FixedLen {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}
	for _, v := range vec.Data {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		body = append(body, buf...)
	}

	codec.Append(vec, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_UINT}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorUintDecode(data []byte) (*AmfVectorUint, int, error) {
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
			if asrt, ok := val.(*AmfVectorUint); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a vector uint")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		if len(data) < 1 {
			return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorUint")
		}
		vec := &AmfVectorUint{}
		vec.FixedLen = data[0] == 0x01
		vec.Data = []uint32{}
		data = data[1:]
		totalConsumed++

		// u29i-value *(uint-value)
		length := U29i >> 1
		if uint32(len(data)) < length * 4 {
			return nil, 0, fmt.Errorf("Vector uint is shorter than expected: %v", data)
		}
		for i := 0; i < int(length); i++ {
			buf := data[i*4:i*4+4]
			vec.Data = append(vec.Data, binary.BigEndian.Uint32(buf))
		}

		codec.Append(vec, COMPLEX_TABLE)
		return vec, totalConsumed, nil
	}
}

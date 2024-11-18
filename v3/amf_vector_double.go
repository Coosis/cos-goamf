package v3

import (
	"fmt"
	"math"
	"encoding/binary"
)

type AmfVectorDouble struct {
	FixedLen bool
	Data []float64
}

func(codec *AmfCodec) AmfVectorDoubleEncode(vec *AmfVectorDouble) ([]byte, error) {
	if id, ok := codec.GetId(vec, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_VECTOR_DOUBLE}, objref...), nil
	}

	length := uint32(len(vec.Data))
	U29V, err := AmfIntEncodePayload(length << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := U29V
	if vec.FixedLen {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}

	for _, v := range vec.Data {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(v))
		body = append(body, buf...)
	}

	codec.Append(vec, COMPLEX_TABLE)
	res := append([]byte{AMF_VECTOR_DOUBLE}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfVectorDoubleDecode(data []byte) (*AmfVectorDouble, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorDouble")
	}
	if data[0] != AMF_VECTOR_DOUBLE {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_VECTOR_DOUBLE, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29V, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29V & 1) == 0 {
		// d-ref
		id := U29V >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.(*AmfVectorDouble); ok {
				return asrt, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a vector double")
		} else {
			return nil, 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		if len(data) < 1 {
			return nil, 0, fmt.Errorf("Not enough data to decode AmfVectorDouble")
		}
		res := &AmfVectorDouble{}
		res.FixedLen = data[0] == 0x01
		res.Data = []float64{}
		data = data[1:]
		totalConsumed++

		// u29v-value *(double-value)
		length := U29V >> 1
		if uint32(len(data)) < length * 8 {
			return nil, 0, fmt.Errorf("Vector double is shorter than expected: %v", data)
		}

		for i := 0; i < int(length); i++ {
			buf := data[i*8:i*8+8]
			res.Data = append(res.Data, math.Float64frombits(binary.BigEndian.Uint64(buf)))
		}

		codec.Append(res, COMPLEX_TABLE)

		return res, totalConsumed, nil
	}
}

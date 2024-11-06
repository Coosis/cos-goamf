package v3

import (
	"fmt"
	"time"
)

func(codec *AmfCodec) AmfDateNow() ([]byte, error) {
	return codec.AmfDateEncode(float64(time.Now().UnixMilli()))
}

func(codec *AmfCodec) AmfDateEncode(value float64) ([]byte, error) {
	if id, ok := codec.GetId(value, COMPLEX_TABLE); ok {
		num, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_DATE}, num...), nil
	}
	res, err := AmfDoubleEncodePayload(value)
	if err != nil {
		return nil, err
	}
	u29d, err := AmfIntEncodePayload(uint32(1))
	if err != nil {
		return nil, err
	}
	res = append(u29d, res...)
	res = append([]byte{AMF_DATE}, res...)
	codec.Append(value, COMPLEX_TABLE)
	return res, nil
}

func(codec *AmfCodec) AmfDateDecode(data []byte) (float64, int, error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfDate")
	}
	if data[0] != AMF_DATE {
		return 0, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_DATE, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	num, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return 0, 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	// raw type
	if (num & 0x01) != 0 {
		if len(data) < 8 {
			return 0, 0, fmt.Errorf("Not enough data to decode raw AmfDate")
		}
		res, cnt, err := AmfDoubleDecodePayload(data)
		if err != nil {
			return 0, 0, err
		}
		totalConsumed += cnt
		codec.Append(res, COMPLEX_TABLE)
		return res, totalConsumed, nil
	}
	// ref type
	id := num >> 1
	if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
		asrt, ok := val.(float64)
		if !ok {
			return 0, 0, fmt.Errorf("Value in table is not a float64")
		}
		return asrt, totalConsumed, nil
	}
	return 0, 0, fmt.Errorf("Reference not found in table")
}

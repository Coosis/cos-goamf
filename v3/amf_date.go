package v3

import (
	"fmt"
	"time"
)

type AmfDate float64

func(codec *AmfCodec) AmfDateNow() ([]byte, error) {
	date := time.Now().UnixMilli()
	return codec.AmfDateEncode(AmfDate(date))
}

func(codec *AmfCodec) AmfDateEncode(value AmfDate) ([]byte, error) {
	if id, ok := codec.GetId(value, COMPLEX_TABLE); ok {
		num, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_DATE}, num...), nil
	}
	dbl := float64(value)
	res, err := AmfDoubleEncodePayload(dbl)
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

func(codec *AmfCodec) AmfDateDecode(data []byte) (AmfDate, int, error) {
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
		codec.Append(AmfDate(res), COMPLEX_TABLE)
		return AmfDate(res), totalConsumed, nil
	}
	// ref type
	id := num >> 1
	if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
		asrt, ok := val.(AmfDate)
		if !ok {
			return 0, 0, fmt.Errorf("Value in table is not a float64")
		}
		return AmfDate(asrt), totalConsumed, nil
	}
	return 0, 0, fmt.Errorf("Reference not found in table")
}

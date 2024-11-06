package v3

import (
	"fmt"
)

func(codec *AmfCodec) AmfXmlEncode(value string) ([]byte, error) {
	if id, ok := codec.GetId(value, COMPLEX_TABLE); ok {
		objref, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_XML}, objref...), nil
	}

	length := len(value)
	U29x, err := AmfIntEncodePayload(uint32(length<<1) | 1)
	if err != nil {
		return nil, err
	}
	codec.Append(value, COMPLEX_TABLE)
	body := append(U29x, []byte(value)...)
	res := append([]byte{AMF_XML}, body...)
	return res, nil
}

func(codec *AmfCodec) AmfXmlDecode(data []byte) (string, int, error) {
	if len(data) < 2 {
		return "", 0, fmt.Errorf("Not enough data to decode AmfXml")
	}
	if data[0] != AMF_XML {
		return "", 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_XML, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29x, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return "", 0, err
	}
	data = data[cnt:]
	totalConsumed += cnt

	if (U29x & 1) == 0 {
		// o-ref
		id := U29x >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if asrt, ok := val.(string); ok {
				return asrt, totalConsumed, nil
			}
			return "", 0, fmt.Errorf("Value in table is not a string")
		} else {
			return "", 0, fmt.Errorf("Reference not found in table: %v", id)
		}
	} else {
		// u29x-value *(utf8-char)
		length := U29x >> 1
		if uint32(len(data)) < length {
			return "", 0, fmt.Errorf("String is shorter than expected: %v", data)
		}
		res := string(data[:length])
		codec.Append(res, COMPLEX_TABLE)
		totalConsumed += int(length)
		return res, totalConsumed, nil
	}
}

package v0

import (
	"fmt"
)

func AmfStringEncode(s string) []byte {
	res := []byte{AMF_STRING}
	utf8 := AmfUTF8Encode(s)
	res = append(res, utf8...)
	return res
}

func AmfStringDecode(data []byte) (string, int, error) {
	if len(data) < 3 {
		return "", 0, fmt.Errorf("Not enough data to decode AmfString")
	}
	if data[0] != AMF_STRING {
		return "", 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected: %T: %v, got: %T: %v",
			AMF_STRING,
			AMF_STRING, 
			data[0],
			data[0],
		)
	}
	data = data[1:]
	str, cnt, err := AmfUTF8Decode(data)
	if err != nil {
		return "", 0, err
	}
	return str, cnt + 1, nil
}

func AmfLongStringEncode(s string) []byte {
	res := []byte{AMF_LONGSTRING}
	utf8 := AmfUTF8LongEncode(s)
	res = append(res, utf8...)
	return res
}

func AmfLongStringDecode(data []byte) (string, int, error) {
	if len(data) < 5 {
		return "", 0, fmt.Errorf("Not enough data to decode AmfLongString")
	}
	if data[0] != AMF_LONGSTRING {
		return "", 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected: %T: %v, got: %T: %v",
			AMF_LONGSTRING,
			AMF_LONGSTRING, 
			data[0],
			data[0],
		)
	}
	data = data[1:]
	str, cnt, err := AmfUTF8LongDecode(data)
	if err != nil {
		return "", 0, err
	}
	return str, cnt + 1, nil
}

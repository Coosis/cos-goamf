package v0

import (
	"fmt"
)

type AmfXmldoc string

func AmfXmldocEncode(doc AmfXmldoc) []byte {
	res := make([]byte, 0)
	res = append(res, AMF_XMLDOC)
	str := string(doc)
	res = append(res, AmfUTF8LongEncode(str)...)
	return res
}

func AmfXmldocDecode(data []byte) (AmfXmldoc, int, error) {
	if len(data) < 5 {
		return "", 0, fmt.Errorf("Not enough data to decode AmfXmldoc")
	}
	if data[0] != AMF_XMLDOC {
		return "", 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected: %T: %v, got: %T: %v",
			AMF_XMLDOC,
			AMF_XMLDOC, 
			data[0],
			data[0],
		)
	}

	data = data[1:]
	str, cnt, err := AmfUTF8LongDecode(data)
	if err != nil {
		return "", 0, err
	}
	return AmfXmldoc(str), cnt + 1, nil
}

package v0

import (
	"fmt"
)

func AmfBooleanEncode(b bool) []byte {
	if b {
		return append([]byte{AMF_BOOLEAN}, byte(1))
	}
	return []byte{AMF_BOOLEAN, byte(0)}
}

func AmfBooleanDecode(data []byte) (bool, int, error) {
	if len(data) < 2 {
		return false, 0, fmt.Errorf("Not enough data to decode AmfBoolean")
	}
	if data[0] != AMF_BOOLEAN {
		return false, 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected: %T: %v, got: %T: %v",
			AMF_BOOLEAN,
			AMF_BOOLEAN, 
			data[0],
			data[0],
		)
	}
	return data[1] != 0, 2, nil
}

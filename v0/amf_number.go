package v0

import (
	"fmt"
	"math"
	"encoding/binary"
)

func AmfNumberEncode(num float64) ([]byte, error) {
	res := []byte{AMF_NUMBER}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, math.Float64bits(num))
	return append(res, buf...), nil
}

func AmfNumberDecode(data []byte) (float64, int, error) {
	if len(data) < 9 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfNumber")
	}
	if data[0] != AMF_NUMBER {
		return 0, 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected: %T: %v, got: %T: %v",
			AMF_NUMBER,
			AMF_NUMBER, 
			data[0],
			data[0],
		)
	}
	data = data[1:]

	if len(data) < 8 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfNumber")
	}
	res := math.Float64frombits(binary.BigEndian.Uint64(data))
	return res, 9, nil
}

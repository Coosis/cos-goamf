package v0

import (
	"fmt"
	"math"
	"encoding/binary"
)

type AmfDate float64

func AmfDateEncode(date AmfDate) []byte {
	buf := make([]byte, 11)
	buf[0] = AMF_DATE
	float64bits := math.Float64bits(float64(date))
	binary.BigEndian.PutUint64(buf[1:9], float64bits)
	binary.BigEndian.PutUint16(buf[9:11], 0)
	return buf
}

func AmfDateDecode(data []byte) (AmfDate, int, error) {
	if len(data) < 11 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfDate")
	}
	if data[0] != AMF_DATE {
		return 0, 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected: %T: %v, got: %T: %v",
			AMF_DATE,
			AMF_DATE, 
			data[0],
			data[0],
		)
	}

	if binary.BigEndian.Uint16(data[9:11]) != 0 {
		return 0, 0, fmt.Errorf("Invalid AmfDate")
	}

	date := math.Float64frombits(binary.BigEndian.Uint64(data[1:]))
	return AmfDate(date), 11, nil
}

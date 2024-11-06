package v3

import(
	"fmt"
	"encoding/binary"
	"math"
)

func AmfDoubleEncode(value float64) ([]byte, error) {
	res, err := AmfDoubleEncodePayload(value)
	if err != nil {
		return nil, err
	}
	return append([]byte{AMF_DOUBLE}, res...), nil
}

func AmfDoubleEncodePayload(value float64) ([]byte, error) {
	if math.IsNaN(value) {
		return nil, fmt.Errorf("Cannot encode NaN as AmfDouble")
	}
	bits := math.Float64bits(value)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, bits)
	return buf, nil
}

func AmfDoubleDecode(data []byte) (float64, int, error) {
	if len(data) < 9 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfDouble")
	}
	if data[0] != AMF_DOUBLE {
		return 0, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_DOUBLE, data[0])
	}
	num, cnt, err := AmfDoubleDecodePayload(data[1:])
	if err != nil {
		return 0, 0, err
	}
	return num, cnt+1, nil
}

func AmfDoubleDecodePayload(data []byte) (float64, int, error) {
	if len(data) < 8 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfDouble")
	}
	bits := binary.BigEndian.Uint64(data[:8])
	return math.Float64frombits(bits), 8, nil
}

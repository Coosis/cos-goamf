package v3

import (
	"math"
	"encoding/binary"
	"testing"
)

func float64bits(value float64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, math.Float64bits(value))
	return buf
}

func genDateTestCases() ([]float64, [][]byte) {
	dates := []float64{
		0.0,
		1.0,
		1.23456789,
		-1.23456789,
		1.0,
		1.0,
		-1.23456789,
		3.14159265,
	}
	flag, _ := AmfIntEncodePayload(1)
	r5, _ := AmfIntEncodePayload(1 << 1 | 0)
	r6, _ := AmfIntEncodePayload(1 << 1 | 0)
	r7, _ := AmfIntEncodePayload(3 << 1 | 0)
	bytes := [][]byte{
		append(flag, float64bits(dates[0])...),
		append(flag, float64bits(dates[1])...),
		append(flag, float64bits(dates[2])...),
		append(flag, float64bits(dates[3])...),
		r5,
		r6,
		r7,
		append(flag, float64bits(dates[7])...),
	}
	for i, b := range bytes {
		b = append([]byte{AMF_DATE}, b...)
		bytes[i] = b
	}
	return dates, bytes
}

func TestAmfDateEncode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genDateTestCases()
	for i, value := range testcases {
		result, err := codec.AmfDateEncode(value)
		if err != nil {
			t.Errorf("AmfDateEncode failed: %v", err)
		}
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfDateEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfDateDecode(t *testing.T) {
	codec := NewAmfCodec()
	expected, testcases := genDateTestCases()
	for i, bytes := range testcases {
		result, _, err := codec.AmfDateDecode(bytes)
		if err != nil {
			t.Errorf("AmfDateDecode failed: %v", err)
		}
		if result != expected[i] {
			t.Errorf("AmfDateDecode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfDateEncodeDecode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genDateTestCases()

	encoded := make([][]byte, len(testcases))

	t.Run("AmfDateEncode", func(t *testing.T) {
		for i, value := range testcases {
			result, err := codec.AmfDateEncode(value)
			if err != nil {
				t.Errorf("AmfDateEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i]) {
				t.Errorf("AmfDateEncode failed: expected %v, got %v", expected[i], result)
			}
			encoded[i] = result
		}
	})

	t.Run("AmfDateDecode", func(t *testing.T) {
		for i, enc := range encoded {
			result, _, err := codec.AmfDateDecode(enc)
			if err != nil {
				t.Errorf("AmfDateDecode failed: %v", err)
			}
			if result != testcases[i] {
				t.Errorf("AmfDateDecode failed: expected %v, got %v", testcases[i], result)
			}
		}
	})
}

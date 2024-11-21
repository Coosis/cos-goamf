package v0

import (
	"testing"
	"math"
	"encoding/binary"
)

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, av := range a {
		if av != b[i] {
			return false
		}
	}
	return true
}

func genNumTestcases() ([]float64, [][]byte) {
	nums := []float64{
		0.0,
		1.0,
		1.23456789,
		-1.23456789,
		1.0,
		1.0,
		-1.23456789,
		3.14159265,
	}

	bytes := [][]byte{}

	for _, num := range nums {
		buf := make([]byte, 8)
		bits := math.Float64bits(num)
		binary.BigEndian.PutUint64(buf, bits)
		bytes = append(bytes, append([]byte{AMF_NUMBER}, buf...))
	}
	return nums, bytes
}

func TestAmfNumberEncode(t *testing.T) {
	testcases, expected := genNumTestcases()
	for i, value := range testcases {
		result, err := AmfNumberEncode(value)
		if err != nil {
			t.Errorf("AmfNumberEncode failed: %v", err)
		}
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfNumberEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfNumberDecode(t *testing.T) {
	testcases, expected := genNumTestcases()
	for i, bytes := range expected {
		result, _, err := AmfNumberDecode(bytes)
		if err != nil {
			t.Errorf("AmfNumberDecode failed: %v", err)
		}
		if result != testcases[i] {
			t.Errorf("AmfNumberDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}

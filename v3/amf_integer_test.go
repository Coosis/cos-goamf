package v3

import (
	"testing"
)

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, value := range a {
		if value != b[i] {
			return false
		}
	}
	return true
}

func genIntTestCases() ([]uint32, [][]byte) {
	nums := []uint32{
		0,
		127,
		128,
		16383,
		16384,
		2097151,
		4194304,
		16,
	}
	bytes := [][]byte{
		{0x00},
		{0x7f},
		{0x81, 0x00},
		{0xff, 0x7f},
		{0x81, 0x80, 0x00},
		{0xff, 0xff, 0x7f},
		{0x81, 0x80, 0x80, 0x00},
		{0x10},
	}
	return nums, bytes
}

func genIntLens() []int {
	return []int{
		1,
		1,
		2,
		2,
		3,
		3,
		4,
		1,
	}
}

func TestAmfIntEncode(t *testing.T) {
	testcases, expected := genIntTestCases()
	for i, value := range testcases {
		result, err := AmfIntEncode(value)
		if err != nil {
			t.Errorf("AmfIntEncode failed: %v", err)
		}
		exp := append([]byte{AMF_INTEGER}, expected[i]...)
		if !bytesEqual(result, exp) {
			t.Errorf("AmfIntEncode %v failed: expected %v, got %v", value, exp, result)
		}
	}
}

func TestAmfIntDecode(t *testing.T) {
	expected, testcases := genIntTestCases()
	lens := genIntLens()
	for i, bytes := range testcases {
		result, cnt, err := AmfIntDecode(append([]byte{AMF_INTEGER}, bytes...))
		if err != nil {
			t.Errorf("AmfIntDecode failed: %v", err)
		}
		if cnt != lens[i]+1 {
			t.Errorf("AmfIntDecode failed: expected %v, got %v", lens[i], cnt)
		}
		if result != expected[i] {
			t.Errorf("AmfIntDecode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfIntDecodePayload(t *testing.T) {
	expected, testcases := genIntTestCases()
	lens := genIntLens()
	for i, bytes := range testcases {
		result, cnt, err := AmfIntDecodePayload(bytes)
		if err != nil {
			t.Errorf("AmfIntDecode failed: %v", err)
		}
		if cnt != lens[i] {
			t.Errorf("AmfIntDecode failed: expected %v, got %v", lens[i], cnt)
		}
		if result != expected[i] {
			t.Errorf("AmfIntDecode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfIntEncodePayload(t *testing.T) {
	testcases, expected := genIntTestCases()
	for i, value := range testcases {
		result, err := AmfIntEncodePayload(value)
		if err != nil {
			t.Errorf("AmfIntEncode failed: %v", err)
		}
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfIntEncode %v failed: expected %v, got %v", value, expected[i], result)
		}
	}
}

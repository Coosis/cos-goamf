package v0

import (
	"testing"
)

func genBoolTestcases() ([]bool, [][]byte) {
	bools := []bool{
		true,
		false,
	}

	bytes := [][]byte{
		{AMF_BOOLEAN, byte(1)},
		{AMF_BOOLEAN, byte(0)},
	}

	return bools, bytes
}

func TestAmfBooleanEncode(t *testing.T) {
	testcases, expected := genBoolTestcases()
	for i, value := range testcases {
		result := AmfBooleanEncode(value)
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfBooleanEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfBooleanDecode(t *testing.T) {
	testcases, expected := genBoolTestcases()
	for i, bytes := range expected {
		result, _, err := AmfBooleanDecode(bytes)
		if err != nil {
			t.Errorf("AmfBooleanDecode failed: %v", err)
		}
		if result != testcases[i] {
			t.Errorf("AmfBooleanDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}

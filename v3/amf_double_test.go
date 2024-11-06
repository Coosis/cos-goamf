package v3

import (
	"testing"
	"math"
)

func genDoubleTestCases() []float64 {
	return []float64{
		0.0,
		1.0,
		1.23456789,
		math.Inf(1),
		math.Inf(-1),
	}
}

func TestAmfDouble(t *testing.T) {
	tests := genDoubleTestCases()
	for _, test := range tests {
		encoded, err := AmfDoubleEncode(test)
		if err != nil {
			t.Errorf("Error encoding AmfDouble: %v", err)
		}
		decoded, cnt, err := AmfDoubleDecode(encoded)
		if err != nil {
			t.Errorf("Error decoding AmfDouble: %v", err)
		}
		if cnt != 9 {
			t.Errorf("Expected length to be 9(header+double), got %v", cnt)
		}
		if test != decoded {
			t.Errorf("Expected %v, got %v", test, decoded)
		}
	}
}

func TestAmfDoublePayload(t *testing.T) {
	tests := genDoubleTestCases()
	for _, test := range tests {
		encoded, err := AmfDoubleEncodePayload(test)
		if err != nil {
			t.Errorf("Error encoding AmfDouble: %v", err)
		}
		decoded, cnt, err := AmfDoubleDecodePayload(encoded)
		if err != nil {
			t.Errorf("Error decoding AmfDouble: %v", err)
		}
		if cnt != 8 {
			t.Errorf("Expected length to be 8(double), got %v", cnt)
		}
		if test != decoded {
			t.Errorf("Expected %v, got %v", test, decoded)
		}
	}
	_, err := AmfDoubleEncodePayload(math.NaN())
	if err == nil {
		t.Errorf("Expected error encoding NaN")
	}

}

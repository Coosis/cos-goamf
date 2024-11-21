package v0

import (
	// "fmt"
	"testing"
	"math"
	"encoding/binary"
)

func genDateTestcases() ([]AmfDate, [][]byte) {
	dates := []AmfDate{
		0,
		1,
		1.1,
		1.2,
	}

	buf1 := make([]byte, 11)
	buf1[0] = AMF_DATE
	bits1 := math.Float64bits(float64(dates[0]))
	binary.BigEndian.PutUint64(buf1[1:9], bits1)
	binary.BigEndian.PutUint16(buf1[9:11], 0)
	// fmt.Println(buf1)

	buf2 := make([]byte, 11)
	buf2[0] = AMF_DATE
	bits2 := math.Float64bits(float64(dates[1]))
	binary.BigEndian.PutUint64(buf2[1:9], bits2)
	binary.BigEndian.PutUint16(buf2[9:11], 0)

	buf3 := make([]byte, 11)
	buf3[0] = AMF_DATE
	bits3 := math.Float64bits(float64(dates[2]))
	binary.BigEndian.PutUint64(buf3[1:9], bits3)
	binary.BigEndian.PutUint16(buf3[9:11], 0)

	buf4 := make([]byte, 11)
	buf4[0] = AMF_DATE
	bits4 := math.Float64bits(float64(dates[3]))
	binary.BigEndian.PutUint64(buf4[1:9], bits4)
	binary.BigEndian.PutUint16(buf4[9:11], 0)

	bytes := [][]byte{
		buf1,
		buf2,
		buf3,
		buf4,
	}

	return dates, bytes
}

func TestAmfDateEncode(t *testing.T) {
	testcases, expected := genDateTestcases()
	for i, value := range testcases {
		result := AmfDateEncode(value)
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfDateEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfDateDecode(t *testing.T) {
	testcases, expected := genDateTestcases()
	for i, bytes := range expected {
		result, _, err := AmfDateDecode(bytes)
		if err != nil {
			t.Errorf("AmfDateDecode failed: %v", err)
		}
		if result != testcases[i] {
			t.Errorf("AmfDateDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}

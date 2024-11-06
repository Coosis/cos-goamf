package v3

import(
	"testing"
)

func genByteArrayTestCases() ([]*[]byte, [][]byte) {
	b1 := []byte{}
	b2 := []byte{0x01, 0x02, 0x03}
	b3 := []byte{0x00, 0x01, 0x02, 0x03, 0x04}

	buf1 := []byte{0x0c}
	int1, _ := AmfIntEncodePayload(uint32(0) << 1 | 1)
	buf1 = append(buf1, int1...)
	buf1 = append(buf1, b1...)
	buf2 := []byte{0x0c}
	int2, _ := AmfIntEncodePayload(uint32(3) << 1 | 1)
	buf2 = append(buf2, int2...)
	buf2 = append(buf2, b2...)
	buf3 := []byte{0x0c}
	int3, _ := AmfIntEncodePayload(uint32(5) << 1 | 1)
	buf3 = append(buf3, int3...)
	buf3 = append(buf3, b3...)
	buf4 := []byte{0x0c}
	int4, _ := AmfIntEncodePayload(uint32(1) << 1)
	buf4 = append(buf4, int4...)

	return []*[]byte{&b1, &b2, &b3, &b2}, [][]byte{buf1, buf2, buf3, buf4}
}

func TestAmfByteArrayEncode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genByteArrayTestCases()
	for i, value := range testcases {
		result, err := codec.AmfByteArrayEncode(value)
		if err != nil {
			t.Errorf("AmfByteArrayEncode failed: %v", err)
		}
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfByteArrayEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

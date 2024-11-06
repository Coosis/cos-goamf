package v3

import (
	"testing"
)

func genStringTestCases() ([]string, [][]byte) {
	strs := []string{
		"",
		"hello",
		"world",
		"hello",
		"helloworld",
		"world",
		"world",
		"This is a new sentence",
	}
	h2, _ := AmfIntEncodePayload(uint32(len(strs[1]))<<1|1)
	h3, _ := AmfIntEncodePayload(uint32(len(strs[2]))<<1|1)
	r4, _ := AmfIntEncodePayload(0 << 1)
	h5, _ := AmfIntEncodePayload(uint32(len(strs[4]))<<1|1)
	r6, _ := AmfIntEncodePayload(1 << 1)
	r7, _ := AmfIntEncodePayload(1 << 1)
	h8, _ := AmfIntEncodePayload(uint32(len(strs[7]))<<1|1)
	return strs,
	[][]byte{
		{0x01},
		append(h2, []byte(strs[1])...),
		append(h3, []byte(strs[2])...),
		r4,
		append(h5, []byte(strs[4])...),
		r6,
		r7,
		append(h8, []byte(strs[7])...),
	}
}

func TestAmfStringEncodeDecodePayload(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genStringTestCases()

	t.Run("AmfStringEncodePayload", func(t *testing.T) {
		for i, value := range testcases {
			result, err := codec.AmfStringEncodePayload(value)
			if err != nil {
				t.Errorf("AmfStringEncodePayload failed: %v", err)
			}
			if !bytesEqual(result, expected[i]) {
				t.Errorf("AmfStringEncodePayload failed: expected %v, got %v", expected[i], result)
			}
		}
	})

	t.Run("AmfStringDecodePayload", func(t *testing.T) {
		for i, bytes := range expected {
			result, _, err := codec.AmfStringDecodePayload(bytes)
			if err != nil {
				t.Errorf("AmfStringDecodePayload failed: %v", err)
			}
			if result != testcases[i] {
				t.Errorf("AmfStringDecodePayload failed: expected %v, got %v", testcases[i], result)
			}
		}
	})
}

func TestAmfStringEncodeDecode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genStringTestCases()

	t.Run("AmfStringEncode", func(t *testing.T) {
		for i, value := range testcases {
			result, err := codec.AmfStringEncode(value)
			if err != nil {
				t.Errorf("AmfStringEncode failed: %v", err)
			}
			if !bytesEqual(result, append([]byte{AMF_STRING}, expected[i]...)) {
				t.Errorf("AmfStringEncode failed: expected %v, got %v", expected[i], result)
			}
		}
	})

	t.Run("AmfStringDecode", func(t *testing.T) {
		for i, bytes := range expected {
			result, _, err := codec.AmfStringDecode(append([]byte{AMF_STRING}, bytes...))
			if err != nil {
				t.Errorf("AmfStringDecode failed: %v", err)
			}
			if result != testcases[i] {
				t.Errorf("AmfStringDecode failed: expected %v, got %v", testcases[i], result)
			}
		}
	})
}

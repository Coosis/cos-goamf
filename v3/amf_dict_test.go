package v3

import (
	"testing"
)

func amfDictEqual(a, b *AmfDict) bool {
	if len(a.EntryKey) != len(b.EntryKey) {
		return false
	}
	if len(a.EntryKeyMarker) != len(b.EntryKeyMarker) {
		return false
	}
	if len(a.EntryValue) != len(b.EntryValue) {
		return false
	}
	for i := range a.EntryKey {
		av := a.EntryKey[i]
		bv := b.EntryKey[i]
		if !valuesEqual(av, bv) {
			return false
		}
	}
	for i := range a.EntryKeyMarker {
		if a.EntryKeyMarker[i] != b.EntryKeyMarker[i] {
			return false
		}
	}
	for i := range a.EntryValue {
		av := a.EntryValue[i]
		bv := b.EntryValue[i]
		if !valuesEqual(av, bv) {
			return false
		}
	}
	for i := range a.EntryValueMarker {
		if a.EntryValueMarker[i] != b.EntryValueMarker[i] {
			return false
		}
	}
	if a.WeakKeys != b.WeakKeys {
		return false
	}
	return true
}

func genDictTestCases() ([]*AmfDict, [][]byte) {
	dict1 := EmptyAmfDict()
	dict1.WeakKeys = false
	dict1.EntryKey = []interface{}{1, 2, 3}
	dict1.EntryKeyMarker = []AmfMarker{AMF_INTEGER, AMF_INTEGER, AMF_INTEGER}
	dict1.EntryValue = []interface{}{1, 2, 3}
	dict1.EntryValueMarker = []AmfMarker{AMF_INTEGER, AMF_INTEGER, AMF_INTEGER}

	dict2 := EmptyAmfDict()
	dict2.WeakKeys = true
	dict2.EntryKey = []interface{}{1, 2, 3}
	dict2.EntryKeyMarker = []AmfMarker{AMF_INTEGER, AMF_INTEGER, AMF_INTEGER}
	dict2.EntryValue = []interface{}{true, false, nil}
	dict2.EntryValueMarker = []AmfMarker{AMF_TRUE, AMF_FALSE, AMF_NULL}

	buf1 := []byte{AMF_DICTIONARY}
	enc1, _ := AmfIntEncodePayload(uint32(3 << 1 | 1))
	enc2, _ := AmfIntEncode(uint32(1))
	enc3, _ := AmfIntEncode(uint32(2))
	enc4, _ := AmfIntEncode(uint32(3))
	buf1 = append(buf1, enc1...)
	buf1 = append(buf1, 0x00)
	buf1 = append(buf1, enc2...)
	buf1 = append(buf1, enc2...)
	buf1 = append(buf1, enc3...)
	buf1 = append(buf1, enc3...)
	buf1 = append(buf1, enc4...)
	buf1 = append(buf1, enc4...)

	buf2 := []byte{AMF_DICTIONARY}
	buf2 = append(buf2, enc1...)
	buf2 = append(buf2, 0x01)
	buf2 = append(buf2, enc2...)
	buf2 = append(buf2, AMF_TRUE)
	buf2 = append(buf2, enc3...)
	buf2 = append(buf2, AMF_FALSE)
	buf2 = append(buf2, enc4...)
	buf2 = append(buf2, AMF_NULL)

	enc5, _ := AmfIntEncodePayload(uint32(0 << 1))
	buf3 := []byte{AMF_DICTIONARY}
	buf3 = append(buf3, enc5...)

	return []*AmfDict{dict1, dict2, dict1}, [][]byte{
		buf1, buf2, buf3,
	}
}

func TestAmfDictEncode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genDictTestCases()
	for i, value := range testcases {
		// t.Logf("TestAmfDictEncode: %v", i)
		result, err := codec.AmfDictEncode(value)
		if err != nil {
			t.Errorf("AmfDictEncode failed: %v", err)
		}
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfDictEncode failed: expected\n%v,\ngot\n%v", expected[i], result)
		}
	}
}

func TestAmfDictDecode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genDictTestCases()
	for i, bytes := range expected {
		// t.Logf("TestAmfDictEncode: %v", i)
		result, _, err := codec.AmfDictDecode(bytes)
		if err != nil {
			t.Errorf("AmfDictDecode failed: %v", err)
		}
		if !amfDictEqual(result, testcases[i]) {
			t.Errorf("AmfDictDecode failed: expected\n%v,\ngot\n%v", testcases[i], result)
		}
	}
}

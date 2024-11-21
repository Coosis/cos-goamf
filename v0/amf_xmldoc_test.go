package v0

import (
	"testing"
	"encoding/binary"
)

func genXmldocTestcases() ([]AmfXmldoc, [][]byte) {
	strs := []string{
		"<some></some>",
		"<some><nested></nested></some>",
		"<some><nested><deep></deep></nested></some>",
		"<b>",
	}
	len1 := make([]byte, 4)
	binary.BigEndian.PutUint32(len1, uint32(len(strs[0])))
	len2 := make([]byte, 4)
	binary.BigEndian.PutUint32(len2, uint32(len(strs[1])))
	len3 := make([]byte, 4)
	binary.BigEndian.PutUint32(len3, uint32(len(strs[2])))
	len4 := make([]byte, 4)
	binary.BigEndian.PutUint32(len4, uint32(len(strs[3])))

	buf1 := append(
		[]byte{AMF_XMLDOC},
		append(len1, []byte(strs[0])...)...
	)

	buf2 := append(
		[]byte{AMF_XMLDOC},
		append(len2, []byte(strs[1])...)...
	)

	buf3 := append(
		[]byte{AMF_XMLDOC},
		append(len3, []byte(strs[2])...)...
	)

	buf4 := append(
		[]byte{AMF_XMLDOC},
		append(len4, []byte(strs[3])...)...
	)

	bytes := [][]byte{
		buf1,
		buf2,
		buf3,
		buf4,
	}

	xmls := make([]AmfXmldoc, 0)
	for _, value := range strs {
		xmls = append(xmls, AmfXmldoc(value))
	}
	return xmls, bytes
}

func TestAmfXmldocEncode(t *testing.T) {
	testcases, expected := genXmldocTestcases()
	for i, value := range testcases {
		result := AmfXmldocEncode(value)
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfXmldocEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfXmldocDecode(t *testing.T) {
	testcases, expected := genXmldocTestcases()
	for i, bytes := range expected {
		result, _, err := AmfXmldocDecode(bytes)
		if err != nil {
			t.Errorf("AmfXmldocDecode failed: %v", err)
		}
		if result != testcases[i] {
			t.Errorf("AmfXmldocDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}

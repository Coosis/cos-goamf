package v3

import (
	"testing"
)

func genXmlTestCases() ([]AmfXml, [][]byte) {
	xmls := []string{
		"<h1></h1>",
		"<h2></h2>",
		"<h1>Hello</h1>",
		"<h2></h2>",
		"<sometag>world</sometag>",
		"<h1>Hello</h1>",
		"<xml><img></img></xml>",
		"<xml><img>someimg.jpg</img></xml>",
		"<xml><img></img></xml>",
	}
	h1, _ := AmfIntEncodePayload(uint32(len(xmls[0]))<<1|1)
	h2, _ := AmfIntEncodePayload(uint32(len(xmls[1]))<<1|1)
	h3, _ := AmfIntEncodePayload(uint32(len(xmls[2]))<<1|1)
	r4, _ := AmfIntEncodePayload(1 << 1)
	h5, _ := AmfIntEncodePayload(uint32(len(xmls[4]))<<1|1)
	r6, _ := AmfIntEncodePayload(2 << 1)
	h7, _ := AmfIntEncodePayload(uint32(len(xmls[6]))<<1|1)
	h8, _ := AmfIntEncodePayload(uint32(len(xmls[7]))<<1|1)
	r9, _ := AmfIntEncodePayload(4 << 1)
	res := make([]AmfXml, len(xmls))
	for i, xml := range xmls {
		res[i] = AmfXml(xml)
	}
	return res, [][]byte{
		append(h1, []byte(xmls[0])...),
		append(h2, []byte(xmls[1])...),
		append(h3, []byte(xmls[2])...),
		r4,
		append(h5, []byte(xmls[4])...),
		r6,
		append(h7, []byte(xmls[6])...),
		append(h8, []byte(xmls[7])...),
		r9,
	}
}

func TestAmfXmlEncodeDecode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genXmlTestCases()

	t.Run("AmfXmlEncode", func(t *testing.T) {
		for i, value := range testcases {
			result, err := codec.AmfXmlEncode(value)
			if err != nil {
				t.Errorf("AmfXmlDocEncode failed: %v", err)
			}
			exp := append([]byte{AMF_XML}, expected[i]...)
			if !bytesEqual(result, exp) {
				// for i, item := range codec.table {
				// 	t.Logf("table[%v]: %v", i, item)
				// }
				t.Errorf("AmfXmlEncode %v failed: expected %v, got %v", value, exp, result)
			}
		}
	})

	t.Run("AmfXmlDecode", func(t *testing.T) {
		for i, bytes := range expected {
			result, _, err := codec.AmfXmlDecode(append([]byte{AMF_XML}, bytes...))
			if err != nil {
				t.Errorf("AmfXmlDecode failed: %v", err)
			}
			if result != testcases[i] {
			// if !bytesEqual(result[1:], testcases[i]) {
				t.Errorf("AmfXmlDecode failed: expected\n%v,\ngot\n%v", testcases[i], result)
			}
		}
	})
}

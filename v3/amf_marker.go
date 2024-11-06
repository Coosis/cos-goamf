package v3

type AmfMarker uint8
const (
	AMF_UNDEFINED = 0x00
	AMF_NULL = 0x01
	AMF_FALSE = 0x02
	AMF_TRUE = 0x03
	AMF_INTEGER = 0x04
	AMF_DOUBLE = 0x05
	AMF_STRING = 0x06
	AMF_XML_DOC = 0x07
	AMF_DATE = 0x08
	AMF_ARRAY = 0x09
	AMF_OBJECT = 0x0a
	AMF_XML = 0x0b
	AMF_BYTE_ARRAY = 0x0c
	AMF_VECTOR_INT = 0x0d
	AMF_VECTOR_UINT = 0x0e
	AMF_VECTOR_DOUBLE = 0x0f
	AMF_VECTOR_OBJECT = 0x10
	AMF_DICTIONARY = 0x11
)

package v3

import (
	"fmt"
	// "math"
	// "encoding/binary"
)

func(codec *AmfCodec) AmfDecode(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("Not enough data to decode Amf")
	}
	switch data[0] {
	case AMF_UNDEFINED:
		return nil, 1, nil
	case AMF_NULL:
		return nil, 1, nil
	case AMF_FALSE:
		return false, 1, nil
	case AMF_TRUE:
		return true, 1, nil
	case AMF_INTEGER:
		return AmfIntDecode(data)
	case AMF_DOUBLE:
		return AmfDoubleDecode(data)
	case AMF_STRING:
		return codec.AmfStringDecode(data)
	case AMF_XML_DOC:
		return codec.AmfXmlDocDecode(data)
	case AMF_DATE:
		return codec.AmfDateDecode(data)
	case AMF_ARRAY:
		return codec.AmfArrayDecode(data)
	case AMF_OBJECT:
		return codec.AmfObjDecode(data)
	case AMF_XML:
		return codec.AmfXmlDecode(data)
	case AMF_BYTE_ARRAY:
		return codec.AmfByteArrayDecode(data)
	case AMF_VECTOR_INT:
		return codec.AmfVectorIntDecode(data)
	case AMF_VECTOR_UINT:
		return codec.AmfVectorUintDecode(data)
	case AMF_VECTOR_DOUBLE:
		return codec.AmfVectorDoubleDecode(data)
	case AMF_VECTOR_OBJECT:
		return codec.AmfVectorObjDecode(data)
	case AMF_DICTIONARY:
		return codec.AmfDictDecode(data)
	}
	return nil, 0, fmt.Errorf("Invalid AMF type: %v", data[0])
}

func(codec *AmfCodec) AmfEncode(value interface{}) ([]byte, error) {
	switch value.(type) {
	// never return undefined
	// 	return AmfUndefined(), nil
	case nil:
		return AmfNull(), nil
	case bool:
		if value.(bool) {
			return AmfBool(true), nil
		}
		return AmfBool(false), nil
	// all kinds of uints
	case uint, uint8, uint16, uint32:
		if val, ok := value.(uint32); ok {
			return AmfIntEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_INTEGER(uint32): %v", value)
	// ints, floats
	case int, int8, int16, int32, int64, float32, float64:
		if val, ok := value.(int); ok {
			return AmfIntEncode(uint32(val))
		}
		if val, ok := value.(float64); ok {
			return AmfDoubleEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_DOUBLE(float64): %v", value)
	case string:
		if val, ok := value.(string); ok {
			return codec.AmfStringEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_STRING(string): %v", value)
	case AmfXmlDoc:
		if val, ok := value.(AmfXmlDoc); ok {
			return codec.AmfXmlDocEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_XML_DOC(string): %v", value)
	case AmfDate:
		if val, ok := value.(AmfDate); ok {
			return codec.AmfDateEncode(val)
		}
	case *AmfArray:
		if val, ok := value.(*AmfArray); ok {
			return codec.AmfArrayEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_ARRAY(*AmfArray): %v", value)
	case *AmfObj:
		if val, ok := value.(*AmfObj); ok {
			return codec.AmfObjEncode(val)
		}
	case AmfXml:
		if val, ok := value.(AmfXml); ok {
			return codec.AmfXmlEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_XML(string): %v", value)
	case []byte:
		if val, ok := value.([]byte); ok {
			return codec.AmfByteArrayEncode(&val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_BYTE_ARRAY([]byte): %v", value)
	case *AmfVectorInt:
		if val, ok := value.(*AmfVectorInt); ok {
			return codec.AmfVectorIntEncode(val)
		}
	case *AmfVectorUint:
		if val, ok := value.(*AmfVectorUint); ok {
			return codec.AmfVectorUintEncode(val)
		}
	case *AmfVectorDouble:
		if val, ok := value.(*AmfVectorDouble); ok {
			return codec.AmfVectorDoubleEncode(val)
		}
	case *AmfVectorObj:
		if val, ok := value.(*AmfVectorObj); ok {
			return codec.AmfVectorObjEncode(val)
		}
	case *AmfDict:
		if val, ok := value.(*AmfDict); ok {
			return codec.AmfDictEncode(val)
		}
	}
	return nil, fmt.Errorf("Invalid value for AMF: %v", value)
}

func AmfUndefined() []byte {
	return []byte{AMF_UNDEFINED}
}

func AmfNull() []byte {
	return []byte{AMF_NULL}
}

func AmfBool(value bool) []byte {
	if value {
		return []byte{AMF_TRUE}
	}
	return []byte{AMF_FALSE}
}


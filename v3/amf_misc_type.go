package v3

import (
	"fmt"
	// "math"
	// "encoding/binary"
)

func(codec *AmfCodec) AmfDecode(data []byte) (interface{}, AmfMarker, int, error) {
	if len(data) == 0 {
		return nil, 0, 0, fmt.Errorf("Not enough data to decode Amf")
	}
	switch data[0] {
	case AMF_UNDEFINED:
		return nil, AMF_UNDEFINED, 1, nil
	case AMF_NULL:
		return nil, AMF_NULL, 1, nil
	case AMF_FALSE:
		return false, AMF_FALSE, 1, nil
	case AMF_TRUE:
		return true, AMF_TRUE, 1, nil
	case AMF_INTEGER:
		res, cnt, err := AmfIntDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return res, AMF_INTEGER, cnt, nil
	case AMF_DOUBLE:
		res, cnt, err := AmfDoubleDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return res, AMF_DOUBLE, cnt, nil
	case AMF_STRING:
		str, cnt, err := codec.AmfStringDecode(data)
		return str, AMF_STRING, cnt, err
	case AMF_XML_DOC:
		str, cnt, err := codec.AmfXmlDocDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return str, AMF_XML_DOC, cnt, nil
	case AMF_DATE:
		res, cnt, err := codec.AmfDateDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return res, AMF_DATE, cnt, nil
	case AMF_ARRAY:
		arr, cnt, err := codec.AmfArrayDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return arr, AMF_ARRAY, cnt, nil
	case AMF_OBJECT:
		obj, cnt, err := codec.AmfObjDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return obj, AMF_OBJECT, cnt, nil
	case AMF_XML:
		str, cnt, err := codec.AmfXmlDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return str, AMF_XML, cnt, nil
	case AMF_BYTE_ARRAY:
		res, cnt, err := codec.AmfByteArrayDecode(data)
		if err != nil {
			return nil, 0, 0, err
		}
		return res, AMF_BYTE_ARRAY, cnt, nil
	}
	return nil, 0, 0, fmt.Errorf("Unknown AMF type: %v", data[0])
}

func(codec *AmfCodec) AmfEncode(value interface{}, marker AmfMarker) ([]byte, error) {
	switch marker {
	case AMF_UNDEFINED:
		return AmfUndefined(), nil
	case AMF_NULL:
		return AmfNull(), nil
	case AMF_FALSE:
		return AmfBool(false), nil
	case AMF_TRUE:
		return AmfBool(true), nil
	case AMF_INTEGER:
		if val, ok := value.(int); ok {
			return AmfIntEncode(uint32(val))
		}
		if val, ok := value.(uint32); ok {
			return AmfIntEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_INTEGER(uint32): %v", value)
	case AMF_DOUBLE:
		if val, ok := value.(float64); ok {
			return AmfDoubleEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_DOUBLE(float64): %v", value)
	case AMF_STRING:
		if val, ok := value.(string); ok {
			return codec.AmfStringEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_STRING(string): %v", value)
	case AMF_XML_DOC:
		if val, ok := value.(string); ok {
			return codec.AmfXmlDocEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_XML_DOC(string): %v", value)
	case AMF_DATE:
		if val, ok := value.(float64); ok {
			return codec.AmfDateEncode(val)
		}
	case AMF_ARRAY:
		if val, ok := value.(*AmfArray); ok {
			return codec.AmfArrayEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_ARRAY(*AmfArray): %v", value)
	case AMF_OBJECT:
		if val, ok := value.(*AmfObj); ok {
			return codec.AmfObjEncode(val)
		}
	case AMF_XML:
		if val, ok := value.(string); ok {
			return codec.AmfXmlEncode(val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_XML(string): %v", value)
	case AMF_BYTE_ARRAY:
		if val, ok := value.([]byte); ok {
			return codec.AmfByteArrayEncode(&val)
		}
		return nil, fmt.Errorf("Invalid value for AMF_BYTE_ARRAY([]byte): %v", value)
	}
	return nil, fmt.Errorf("Unknown marker: %v", marker)
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


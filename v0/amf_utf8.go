package v0

import (
	"fmt"
	"encoding/binary"
)

// If the string is longer than this, it will be truncated to `AMF_STRING_MAXLEN`
func AmfUTF8Encode(s string) []byte {
	if len(s) > AMF_STRING_MAXLEN {
		s = s[:AMF_STRING_MAXLEN]
	}
	buf := []byte(s)
	res := append([]byte{0, 0}, buf...)
	binary.BigEndian.PutUint16(res[:2], uint16(len(buf)))
	return res
}

func AmfUTF8Decode(data []byte) (string, int, error) {
	if len(data) < 2 {
		return "", 0, fmt.Errorf("Not enough data for a valid AmfUTF8 string")
	}
	bytelen := binary.BigEndian.Uint16(data[:2])
	if len(data) < int(bytelen) + 2 {
		return "", 0, fmt.Errorf("Not enough UTF8 to decode into string. Maybe the string is corrupted")
	}
	length := int(bytelen)
	return string(data[2:length+2]), length + 2, nil
}


func AmfUTF8LongEncode(s string) []byte {
	buf := []byte(s)
	res := append([]byte{0, 0, 0, 0}, buf...)
	binary.BigEndian.PutUint32(res[:4], uint32(len(buf)))
	return res
}

func AmfUTF8LongDecode(data []byte) (string, int, error) {
	if len(data) < 4 {
		return "", 0, fmt.Errorf("Not enough data for a valid AmfUTF8Long string")
	}
	bytelen := binary.BigEndian.Uint32(data[:4])
	if len(data) < int(bytelen) + 4 {
		return "", 0, fmt.Errorf("Not enough UTF8 to decode into string. Maybe the string is corrupted")
	}
	length := int(bytelen)
	return string(data[4:length+4]), length + 4, nil
}

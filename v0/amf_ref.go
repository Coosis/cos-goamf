package v0

import (
	"fmt"
	"encoding/binary"
)

func AmfRefEncode(ref uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, ref)
	return append([]byte{AMF_REFERENCE}, buf...)
}

func AmfRefDecode(data []byte) (uint16, int, error) {
	if len(data) < 3 {
		return 0, 0, fmt.Errorf("Not enough data to decode AmfRef")
	}
	if data[0] != AMF_REFERENCE {
		return 0, 0, fmt.Errorf(
			"AMF_TYPE_MISMATCH, expected marker: %v (AMF_REFERENCE), got marker: %v",
			AMF_REFERENCE, 
			data[0],
		)
	}

	ref := binary.BigEndian.Uint16(data[1:3])
	return ref, 3, nil
}

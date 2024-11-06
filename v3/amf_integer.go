package v3

import (
	"fmt"
)

const (
	HIGH_BIT_MASK = 0x80
	LOW_BITS_MASK = 0x7f
)

type AmfInt uint32

func AmfIntDecode(bytes []byte) (uint32, int, error) {
	if len(bytes) < 2 {
		return 0, 0, fmt.Errorf("AMF_RANGE_EXCEPTION, bytes: %v", bytes)
	}
	if bytes[0] != AMF_INTEGER {
		return 0, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_INTEGER, bytes[0])
	}
	res, cnt, err := AmfIntDecodePayload(bytes[1:])
	if err != nil {
		return 0, 0, err
	}
	return res, cnt+1, nil
}

func AmfIntDecodePayload(bytes []byte) (uint32, int, error) {
	if len(bytes) == 0 {
		return 0, 0, fmt.Errorf("Empty data passed to AmfIntDecodePayload")
	}
	var result uint32
	i := 0
	for i = 0; i < len(bytes); i++ {
		if i == 3 {
			result = (result << 8) | uint32(bytes[i])
			break
		}
		result = (result << 7) | uint32(bytes[i] & LOW_BITS_MASK)
		if bytes[i] & HIGH_BIT_MASK == 0 {
			break
		}
	}
	return result, i+1, nil
}

func AmfIntEncode(value uint32) ([]byte, error) {
	payload, err := AmfIntEncodePayload(value)
	if err != nil {
		return nil, err
	}
	return append([]byte{AMF_INTEGER}, payload...), nil
}

func AmfIntEncodePayload(value uint32) ([]byte, error) {
	if 0x00000000 <= value && value <= 0x0000007f {
		return []byte{
			byte(value),
		}, nil
	} else if 0x00000080 <= value && value <= 0x00003fff {
		return []byte{
			byte((value >> 7) & 0x7f | HIGH_BIT_MASK),
			byte(value & 0x7f),
		}, nil
	} else if 0x00004000 <= value && value <= 0x001fffff {
		return []byte{
			byte((value >> 14) & 0x7f | HIGH_BIT_MASK),
			byte((value >> 7) & 0x7f | HIGH_BIT_MASK),
			byte(value & 0x7f),
		}, nil
	} else if 0x00200000 <= value && value <= 0x3fffffff {
		return []byte{
			byte((value >> 22) & 0x7f | HIGH_BIT_MASK),
			byte((value >> 15) & 0x7f | HIGH_BIT_MASK),
			byte((value >> 8) & 0x7f | HIGH_BIT_MASK),
			byte(value & 0xff),
		}, nil
	} else {
		return nil, fmt.Errorf("AMF_RANGE_EXCEPTION, value: %v", value)
	}
}

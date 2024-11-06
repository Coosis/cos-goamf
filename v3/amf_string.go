package v3

import(
	"fmt"
)

const (
	AMF_STRING_HEADER_MAXLEN = 0xf0 << 24
)

func AmfStringHeader(len uint32) ([]byte, error) {
	if len & AMF_STRING_HEADER_MAXLEN != 0 {
		return nil, fmt.Errorf("Length for string is greater than 2^28-1, value: %v", len)
	}
	len = len << 1 | 1
	encoded, err := AmfIntEncodePayload(len)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func AmfStringRef(id uint32) ([]byte, error) {
	encoded, err := AmfIntEncodePayload(id << 1)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func(codec *AmfCodec) AmfStringEncode(value string) ([]byte, error) {
	res, err := codec.AmfStringEncodePayload(value)
	if err != nil {
		return nil, err
	}
	return append([]byte{AMF_STRING}, res...), nil
}

func(codec *AmfCodec) AmfStringEncodePayload(value string) ([]byte, error) {
	// empty string
	if value == "" {
		return []byte{0x01}, nil
	}

	// string reference
	if id, ok := codec.GetId(value, STRING_TABLE); ok {
		return AmfStringRef(id)
	}

	// string literal
	header, err := AmfStringHeader(uint32(len(value)))
	if err != nil {
		return nil, err
	}

	// add to table
	codec.Append(value, STRING_TABLE)
	
	return append(header, []byte(value)...), nil
}

func(codec *AmfCodec) AmfStringDecode(data []byte) (string, int, error) {
	if len(data) == 0 {
		return "", 0, fmt.Errorf("Empty byte array")
	}
	if data[0] != AMF_STRING {
		return "", 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_STRING, data[0])
	}
	data = data[1:]

	str, cnt, err := codec.AmfStringDecodePayload(data)
	if err != nil {
		return "", 0, err
	}
	return str, cnt + 1, nil
}

func(codec *AmfCodec) AmfStringDecodePayload(bytes []byte) (string, int, error) {
	if len(bytes) == 0 {
		return "", 0, fmt.Errorf("Empty byte array")
	}

	if bytes[0] == 0x01 {
		return "", 1, nil
	}
	
	totalcnt := 0
	header, cnt, err := AmfIntDecodePayload(bytes)
	if err != nil {
		return "", 0, err
	}
	bytes = bytes[cnt:]
	totalcnt += cnt
	// ref
	if header & 1 == 0 {
		id := header >> 1
		if str, ok := codec.Get(id, STRING_TABLE); ok {
			if str, ok := str.(string); ok {
				return str, cnt, nil
			}
			return "", 0, fmt.Errorf("Value in table is not a string")
		}
		return "", 0, fmt.Errorf("String reference not found: %v", id)
	}

	// literal
	length := header >> 1
	if uint32(len(bytes)) < length {
		return "", 0, fmt.Errorf("String is shorter than expected: %v", bytes)
	}
	res := string(bytes[:length])
	codec.Append(res, STRING_TABLE)
	totalcnt += int(length)
	return res, totalcnt, nil
}

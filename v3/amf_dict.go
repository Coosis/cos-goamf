package v3

import (
	"fmt"
)

type AmfDict struct {
	WeakKeys bool
	EntryKey []interface{}
	EntryValue []interface{}
}

func EmptyAmfDict() *AmfDict {
	return &AmfDict{
		WeakKeys: false,
		EntryKey: make([]interface{}, 0),
		EntryValue: make([]interface{}, 0),
	}
}

func(codec *AmfCodec) AmfDictEncode(dict *AmfDict) ([]byte, error) {
	if dict == nil {
		return nil, fmt.Errorf("Empty AmfDict")
	}
	if id, ok := codec.GetId(dict, COMPLEX_TABLE); ok {
		num, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_DICTIONARY}, num...), nil
	}

	// raw
	numEntry := uint32(len(dict.EntryKey))
	num, err := AmfIntEncodePayload(numEntry << 1 | 1)
	if err != nil {
		return nil, err
	}
	body := num
	if dict.WeakKeys {
		body = append(body, 0x01)
	} else {
		body = append(body, 0x00)
	}

	for i := 0; i < int(numEntry); i++ {
		key, err := codec.AmfEncode(dict.EntryKey[i])
		if err != nil {
			return nil, err
		}
		body = append(body, key...)
		// fmt.Println("key", key)

		value, err := codec.AmfEncode(dict.EntryValue[i])
		if err != nil {
			return nil, err
		}
		body = append(body, value...)
	}
	codec.Append(dict, COMPLEX_TABLE)
	return append([]byte{AMF_DICTIONARY}, body...), nil
}

func(codec *AmfCodec) AmfDictDecode(data []byte) (*AmfDict, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("Empty AmfDict")
	}
	if data[0] != AMF_DICTIONARY {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_DICTIONARY, data[0])
	}
	data = data[1:]
	totalConsumed := 1

	U29, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	totalConsumed += cnt
	data = data[cnt:]

	if (U29 & 0x01) == 0 {
		// ref
		id := U29 >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if dict, ok := val.(*AmfDict); ok {
				return dict, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a dictionary: %v", val)
		}
		return nil, 0, fmt.Errorf("Invalid reference, id: %v", id)
	}

	// raw
	dict := EmptyAmfDict()
	numEntry := U29 >> 1
	dict.WeakKeys = data[0] == 0x01
	data = data[1:]
	totalConsumed++

	for i := 0; i < int(numEntry); i++ {
		key, cnt, err := codec.AmfDecode(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt
		dict.EntryKey = append(dict.EntryKey, key)

		value, cnt, err := codec.AmfDecode(data)
		if err != nil {
			return nil, 0, err
		}
		data = data[cnt:]
		totalConsumed += cnt
		dict.EntryValue = append(dict.EntryValue, value)
	}
	codec.Append(dict, COMPLEX_TABLE)
	return dict, totalConsumed, nil
}

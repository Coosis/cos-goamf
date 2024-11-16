package v3

import (
	"fmt"
)

// obj = obj_marker (
// 	U290-ref |
// 	(
// 		U290-traits-ext class-name *(U8)
// 	) |
// 	U290-traits-ref |
// 	(
// 		U290-traits class-name *(U8vr)
// 	)
// ) *(value-type) *(dynamic-member)
// traits: Anonymous, Typed, Externalizable, Dynamic

// TODO: add type for member and dynamic-member

type AmfObjMember struct {
	Key string
	Value interface{}
	Marker AmfMarker
}

type AmfObj struct {
	ClassName string
	Member []AmfObjMember
	DynMembers []AmfObjMember
	ExtTraits []byte
}

func EmptyAmfObj() *AmfObj {
	return &AmfObj{
		ClassName: "",
		Member: make([]AmfObjMember, 0),
		DynMembers: make([]AmfObjMember, 0),
		// Member: make(map[string]interface{}, 0),
		// MemberMarker: make(map[string]AmfMarker, 0),
		// DynMembers: make(map[string]interface{}, 0),
		// DynMemberMarker: make(map[string]AmfMarker, 0),
		ExtTraits: make([]uint8, 0),
	}
}

func(obj *AmfObj) AppendMember(mem AmfObjMember) {
	obj.Member = append(obj.Member, mem)
}

func(obj *AmfObj) AppendDynMember(mem AmfObjMember) {
	obj.DynMembers = append(obj.DynMembers, mem)
}

// used when decoding traits-ext
// given data without obj marker, U29o header, class-name, return traits-ext(U8s)
type AmfTraitExtHandler func(codec *AmfCodec, data []byte) ([]byte, int, error)

func(codec *AmfCodec) AmfObjEncode(obj *AmfObj) ([]byte, error) {
	if obj == nil {
		return nil, fmt.Errorf("Empty AmfObj")
	}
	// U29o-ref
	if id, ok := codec.GetId(obj, COMPLEX_TABLE); ok {
		num, err := AmfIntEncodePayload(id << 1)
		if err != nil {
			return nil, err
		}
		return append([]byte{AMF_OBJECT}, num...), nil
	}

	// U29-traits-ext
	if len(obj.ExtTraits) != 0 {
		num, err := AmfIntEncodePayload(uint32(7))
		if err != nil {
			return nil, err
		}
		traits := NewAmfTraitSet()
		traits.ClassName = obj.ClassName
		traits.IsExternalizable = true
		codec.Append(traits, TRAIT_TABLE)
		res := append([]byte{AMF_OBJECT}, num...)
		name, err := codec.AmfStringEncodePayload(obj.ClassName)
		if err != nil {
			return nil, err
		}
		res = append(res, name...)
		res = append(res, obj.ExtTraits...)
		codec.Append(obj, COMPLEX_TABLE)
		return res, nil
	}

	traits := NewAmfTraitSet()
	traits.ClassName = obj.ClassName
	traits.Traits = make([]string, 0)
	for member := range obj.Member {
		// traits.Traits = append(traits.Traits, key)
		traits.Traits = append(traits.Traits, obj.Member[member].Key)
	}
	traits.IsDynamic = len(obj.DynMembers) != 0
	traits.IsExternalizable = false
	codec.Append(traits, TRAIT_TABLE)
	if id, ok := codec.GetId(traits, TRAIT_TABLE); ok {
		// U29-traits-ref
		num, err := AmfIntEncodePayload(id << 2 | 1)
		if err != nil {
			return nil, err
		}
		res := append([]byte{AMF_OBJECT}, num...)
		codec.Append(obj, COMPLEX_TABLE)
		return res, nil
	}

	// U29-traits
	size := uint32(len(traits.Traits)) << 4
	if traits.IsDynamic {
		size |= 0x08
	}
	size |= 0x03
	num, err := AmfIntEncodePayload(size)
	if err != nil {
		return nil, err
	}
	res := []byte{AMF_OBJECT}
	res = append(res, num...)
	classname, err := codec.AmfStringEncodePayload(traits.ClassName)
	if err != nil {
		return nil, err
	}
	res = append(res, classname...)
	for _, trait := range traits.Traits {
		name, err := codec.AmfStringEncodePayload(trait)
		if err != nil {
			return nil, err
		}
		res = append(res, name...)
	}
	for _, mem := range obj.Member {
		encoded, err := codec.AmfEncode(mem.Value, mem.Marker)
		// fmt.Println("before: ", mem.Value)
		// fmt.Println("encoded: ", encoded)
		// encoded, err := codec.AmfEncode(mem, obj.MemberMarker[trait])
		if err != nil {
			return nil, err
		}
		res = append(res, encoded...)
	}

	if traits.IsDynamic {
		for _, mem := range obj.DynMembers {
			name, err := codec.AmfStringEncodePayload(mem.Key)
			if err != nil {
				return nil, err
			}
			res = append(res, name...)
			// encoded, err := codec.AmfEncode(mem, obj.DynMemberMarker[key])
			encoded, err := codec.AmfEncode(mem.Value, mem.Marker)
			if err != nil {
				return nil, err
			}
			res = append(res, encoded...)
		}
		res = append(res, 0x01)
	}

	codec.Append(obj, COMPLEX_TABLE)
	return res, nil
}

func(codec *AmfCodec) AmfObjDecode(data []byte) (*AmfObj, int, error) {
	totalConsumed := 0
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("Not enough data to decode AmfObj")
	}
	if data[0] != AMF_OBJECT {
		return nil, 0, fmt.Errorf("AMF_TYPE_MISMATCH, expected: %v, got: %v", AMF_OBJECT, data[0])
	}

	// consume marker
	data = data[1:]
	totalConsumed++

	num, cnt, err := AmfIntDecodePayload(data)
	if err != nil {
		return nil, 0, err
	}
	totalConsumed += cnt
	data = data[cnt:]

	if num & 0x01 == 0 {
		// ref
		// complex-object-id
		id := num >> 1
		if val, ok := codec.Get(id, COMPLEX_TABLE); ok {
			if obj, ok := val.(*AmfObj); ok {
				return obj, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not an AmfObj")
		}
		return nil, 0, fmt.Errorf("Object reference not found: %v", id)
	}
	if num & 0x03 == 1 {
		// traits-ref
		// id is trait-ref-id
		id := num >> 2
		if val, ok := codec.Get(id, TRAIT_TABLE); ok {
			if traits, ok := val.(AmfTraitSet); ok {
				obj := EmptyAmfObj()
				obj.ClassName = traits.ClassName
				for _, trait := range traits.Traits {
					val, marker, cnt, err := codec.AmfDecode(data)
					if err != nil {
						return nil, 0, err
					}
					totalConsumed += cnt
					data = data[cnt:]
					// obj.Member[trait] = val
					// obj.MemberMarker[trait] = marker
					obj.AppendMember(AmfObjMember{trait, val, marker})
				}
				codec.Append(obj, COMPLEX_TABLE)
				return obj, totalConsumed, nil
			}
			return nil, 0, fmt.Errorf("Value in table is not a traitset")
		}
	}
	if num & 0x07 == 7 {
		// traits-ext
		trait := NewAmfTraitSet()
		name, cnt, err := codec.AmfStringDecodePayload(data)
		if err != nil {
			return nil, 0, err
		}
		trait.ClassName = name
		trait.IsExternalizable = true
		codec.Append(trait, TRAIT_TABLE)
		totalConsumed += cnt
		data = data[cnt:]
		res, cnt, err := []byte{}, 0, nil
		if codec.externalTraitHandler == nil {
			// default behavior: take-all
			res, cnt, err = data, len(data), nil
		} else {
			res, cnt, err = codec.externalTraitHandler(codec, data)
		}
		if err != nil {
			return nil, 0, err
		}
		totalConsumed += cnt
		obj := EmptyAmfObj()
		obj.ClassName = name
		obj.ExtTraits = res
		codec.Append(obj, COMPLEX_TABLE)
		return obj, totalConsumed, nil
	}
	if num & 0x07 == 3 {
		// traits
		classname, cnt, err := codec.AmfStringDecodePayload(data)
		if err != nil {
			return nil, 0, err
		}
		totalConsumed += cnt
		data = data[cnt:]
		numTraits := num >> 4
		if len(data) < int(numTraits) {
			return nil, 0, fmt.Errorf("Not enough data to decode traits")
		}
		traits := make([]string, numTraits)
		for i := uint32(0); i < numTraits; i++ {
			str, cnt, err := codec.AmfStringDecodePayload(data)
			if err != nil {
				return nil, 0, err
			}
			traits[i] = str
			data = data[cnt:]
			totalConsumed += cnt
		}

		obj := EmptyAmfObj()
		obj.ClassName = classname
		dynamic := ((num >> 3) & 1) == 1
		for _, trait := range traits {
			val, marker, cnt, err := codec.AmfDecode(data)
			if err != nil {
				return nil, 0, err
			}
			totalConsumed += cnt
			data = data[cnt:]
			// obj.Member[trait] = val
			// obj.MemberMarker[trait] = marker
			obj.AppendMember(AmfObjMember{trait, val, marker})
		}
		if dynamic {
			for {
				if len(data) == 0 {
					break
				}
				key, cnt, err := codec.AmfStringDecodePayload(data)
				if err != nil {
					return nil, 0, err
				}
				totalConsumed += cnt
				data = data[cnt:]
				if key == "" {
					break
				}
				if len(data) == 0 {
					break
				}
				amfval, marker, cnt, err := codec.AmfDecode(data)
				if err != nil {
					return nil, 0, err
				}
				if cnt == 0 {
					return nil, 0, fmt.Errorf("AmfDecode fails to consume data")
				}
				totalConsumed += cnt
				data = data[cnt:]
				// obj.DynMembers[key] = amfval
				// obj.DynMemberMarker[key] = marker
				obj.AppendDynMember(AmfObjMember{key, amfval, marker})
			}
		}
		codec.Append(obj, COMPLEX_TABLE)
		return obj, totalConsumed, nil
	}
	return nil, 0, fmt.Errorf("Invalid Object: %v", num)
}

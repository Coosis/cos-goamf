package v3

import (
	// "fmt"
)

type AmfCodecTable int
const (
	COMPLEX_TABLE = iota
	STRING_TABLE
	TRAIT_TABLE
)

type AmfCodec struct {
	table map[interface{}]uint32
	revTable map[uint32]interface{}

	strTable map[string]uint32
	strRevTable map[uint32]string

	traitTable map[string]uint32
	traitRevTable map[uint32]AmfTraitSet

	externalTraitHandler AmfTraitExtHandler
}

func NewAmfCodec() *AmfCodec {
	return &AmfCodec{
		table: make(map[interface{}]uint32),
		revTable: make(map[uint32]interface{}),
		strTable: make(map[string]uint32),
		strRevTable: make(map[uint32]string),
		traitTable: make(map[string]uint32),
		traitRevTable: make(map[uint32]AmfTraitSet),
	}
}

func codecHash(value interface{}) interface{} {
	switch value.(type) {
	// case *AmfObj:
	// 	fmt.Println("Hashed!")
	// 	return value.(*AmfObj).Hash()
	default:
		return value
	}
}

func(codec *AmfCodec) Append(value interface{}, kind AmfCodecTable) bool {
	switch kind {
	case COMPLEX_TABLE:
		tableLen := uint32(len(codec.table))
		valueHash := codecHash(value)
		codec.table[valueHash] = tableLen
		codec.revTable[tableLen] = value
	case STRING_TABLE:
		value, ok := value.(string)
		if !ok {
			return false
		}
		tableLen := uint32(len(codec.strTable))
		codec.strTable[value] = tableLen
		codec.strRevTable[tableLen] = value
	case TRAIT_TABLE:
		value, ok := value.(AmfTraitSet)
		hash := value.Hash()
		if !ok {
			return false
		}
		tableLen := uint32(len(codec.traitTable))
		codec.traitTable[hash] = tableLen
		codec.traitRevTable[tableLen] = value
	}
	return true
}

func(codec *AmfCodec) Get(id uint32, kind AmfCodecTable) (interface{}, bool) {
	switch kind {
	case COMPLEX_TABLE:
		if value, ok := codec.revTable[id]; ok {
			return value, true
		}
	case STRING_TABLE:
		if value, ok := codec.strRevTable[id]; ok {
			return value, true
		}
	case TRAIT_TABLE:
		if value, ok := codec.traitRevTable[id]; ok {
			return value, true
		}
	}
	return nil, false
}

func (codec *AmfCodec) GetId(value interface{}, kind AmfCodecTable) (uint32, bool) {
	switch kind {
	case COMPLEX_TABLE:
		value = codecHash(value)
		if id, ok := codec.table[value]; ok {
			return id, ok
		}
	case STRING_TABLE:
		if value, ok := value.(string); ok {
			if id, ok := codec.strTable[value]; ok {
				return id, ok
			}
		}
	case TRAIT_TABLE:
		if value, ok := value.(AmfTraitSet); ok {
			hash := value.Hash()
			if id, ok := codec.traitTable[hash]; ok {
				return id, ok
			}
		}
	}
	return 0, false
}

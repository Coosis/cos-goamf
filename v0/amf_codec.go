package v0

import (
	"fmt"
	"sync"
)

const (
	AMF_STRING_MAXLEN = 65535
)

type AmfCodec struct {
	table map[uint16]interface{}
	revTable map[interface{}]uint16

	mu sync.RWMutex
}

func NewAmfCodec() *AmfCodec {
	return &AmfCodec{
		table: make(map[uint16]interface{}),
		revTable: make(map[interface{}]uint16),
	}
}

func(c *AmfCodec) GetId(val interface{}) (uint16, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if id, ok := c.revTable[val]; ok {
		return id, nil
	}
	return 0, fmt.Errorf("Value not found in table")
}

func(c *AmfCodec) GetObj(id uint16) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.table[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Value not found in table")
}

func(c *AmfCodec) Set(val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.revTable[val]; ok {
		return
	}
	id := uint16(len(c.table))
	c.table[id] = val
	c.revTable[val] = id
}

func(c *AmfCodec) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.table = make(map[uint16]interface{})
	c.revTable = make(map[interface{}]uint16)
}

// Basically asserts the type of val and returns it as a float64.
// Only pass in uint*, int*.
func toFloat64(val interface{}) (float64, error) {
	switch val.(type) {
	case int:
		return float64(val.(int)), nil
	case int8:
		return float64(val.(int8)), nil
	case int16:
		return float64(val.(int16)), nil
	case int32:
		return float64(val.(int32)), nil
	case int64:
		return float64(val.(int64)), nil
	case uint:
		return float64(val.(uint)), nil
	case uint8:
		return float64(val.(uint8)), nil
	case uint16:
		return float64(val.(uint16)), nil
	case uint32:
		return float64(val.(uint32)), nil
	case uint64:
		return float64(val.(uint64)), nil
	case float64:
		return val.(float64), nil
	default:
		return 0, fmt.Errorf("Unsupported type")
	}
}

// If a type is implemented because primitive types are not 
// sufficient, then that type would be implemented as a struct. 
// When passing such a type, they should be passed as a pointer.
// Specifically, the following types needs to be passed as pointers:
// 1. AmfObj
// 2. AmfECMA
// 3. AmfArray
// Passing primitive types will result in the following encoding:
// uint*, int* -> float64 -> AmfNumber
// float64     -> AmfNumber
// bool        -> AmfBoolean
// string      -> AmfStr | AmfLongStr
// nil         -> AmfNull
// Some types can never come out of Encode():
//  x -> AmfMovieclip (not used when encoding)
//  x -> AmfUndefined (not used when encoding)
//  x -> AmfUnsupported (not used when encoding)
//  x -> AmfRecordset (not used when encoding)
func(c *AmfCodec) Encode(val interface{}) ([]byte, error) {
	switch val.(type) {
	case string:
		strlen := len(val.(string))
		if strlen > AMF_STRING_MAXLEN {
			return AmfLongStringEncode(val.(string)), nil
		}
		return AmfStringEncode(val.(string)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		val, err := toFloat64(val)
		if err != nil {
			return nil, err
		}
		return AmfNumberEncode(val)
	case float64:
		return AmfNumberEncode(val.(float64))
	case bool:
		return AmfBooleanEncode(val.(bool)), nil
	case nil:
		return AmfNull(), nil

	// reference type will be used when encoding: 
	// 1. obj 2. array 3. ecma-array

	// Obj type
	case *AmfObj:
		v := val.(*AmfObj)
		if v == nil {
			return AmfNull(), nil
		}
		return c.AmfObjEncode(v)
	// Ecma type
	case *AmfECMA:
		v := val.(*AmfECMA)
		if v == nil {
			return AmfNull(), nil
		}
		return c.AmfECMAEncode(v)
	// Strict Array type
	case *AmfArray:
		v := val.(*AmfArray)
		if v == nil {
			return AmfNull(), nil
		}
		return c.AmfArrayEncode(v)
	case AmfXmldoc:
		return AmfXmldocEncode(val.(AmfXmldoc)), nil
	case AmfDate:
		return AmfDateEncode(val.(AmfDate)), nil
	}
	return nil, fmt.Errorf("Unsupported type")
}

// Notice both Null and Undefined are decoded as `nil` in golang.
// Also, when encountered with the following types:
// 1. AMF_MOVIECLIP
// 2. AMF_UNDEFINED
// 3. AMF_UNSUP
// 4. AMF_RECORDSET
// Decode will return a `Unsupported type` error. Consider 
// handling these types gracefully, instead of panicking.
func(c *AmfCodec) Decode(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("Empty data")
	}
	unsupported := fmt.Errorf("Unsupported type")
	switch data[0] {
	case AMF_NUMBER:
		return AmfNumberDecode(data)
	case AMF_BOOLEAN:
		return AmfBooleanDecode(data)
	case AMF_STRING:
		return AmfStringDecode(data)
	case AMF_OBJECT:
		return c.AmfObjDecode(data)
	case AMF_MOVIECLIP:
		return nil, 1, unsupported
	case AMF_NULL:
		return nil, 1, nil
	case AMF_UNDEFINED:
		return nil, 1, nil
	case AMF_REFERENCE:
		id, cnt, err := AmfRefDecode(data)
		if err != nil {
			return nil, 0, err
		}
		if obj, err := c.GetObj(id); err == nil {
			return obj, cnt, nil
		}
		return nil, 0, fmt.Errorf("Value not found in table")
	case AMF_ECMA:
		return c.AmfECMADecode(data)
	case AMF_OBJEND:
		return AmfObjEnd(), 3, nil
	case AMF_STRICTARR:
		return c.AmfArrayDecode(data)
	case AMF_DATE:
		return AmfDateDecode(data)
	case AMF_LONGSTRING:
		return AmfLongStringDecode(data)
	case AMF_UNSUP:
		return nil, 1, unsupported
	case AMF_RECORDSET:
		return nil, 1, unsupported
	case AMF_XMLDOC:
		return AmfXmldocDecode(data)
	case AMF_TYPEDOBJ:
		return c.AmfObjDecode(data)
	}
	return nil, 0, fmt.Errorf("Unsupported type: %v", data[0])
}

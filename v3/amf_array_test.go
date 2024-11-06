package v3

import (
	"fmt"
	"reflect"
	"testing"
)

func (codec *AmfCodec) LogStringTable(tag string) {
    fmt.Printf("String Table (%s):\n", tag)
    for i, s := range codec.strTable {
        fmt.Printf("%v: %v\n", i, s)
    }
}

func amfArrayEqual(a, b *AmfArray) bool {
    if len(a.Dense) != len(b.Dense) {
        return false
    }
    if len(a.Assoc) != len(b.Assoc) {
        return false
    }
    for i := range a.Dense {
        av := a.Dense[i]
        bv := b.Dense[i]
        if !valuesEqual(av, bv) {
            return false
        }
    }
	for i := range len(a.AssocKeys) {
		if a.AssocKeys[i] != b.AssocKeys[i] {
			return false
		}
		keya := a.Assoc[a.AssocKeys[i]]
		keyb := b.Assoc[b.AssocKeys[i]]
		if !valuesEqual(keya, keyb) {
			return false
		}
		typea := a.AssocMarker[a.AssocKeys[i]]
		typeb := b.AssocMarker[b.AssocKeys[i]]
		if typea != typeb {
			return false
		}
	}
	return true
}

func valuesEqual(a, b interface{}) bool {
	if a==nil && b==nil {
		return true
	}
	if a==nil || b==nil {
		return false
	}
	if numa, ok := a.(int); ok {
		a = uint32(numa)
	}
	if numb, ok := b.(int); ok {
		b = uint32(numb)
	}
    switch av := a.(type) {
    case string:
        bv, ok := b.(string)
        return ok && av == bv
    case bool:
        bv, ok := b.(bool)
        return ok && av == bv
    case uint32:
        bv, ok := b.(uint32)
        return ok && av == bv
    case float64:
        bv, ok := b.(float64)
        return ok && av == bv
    default:
        return reflect.DeepEqual(a, b)
    }
}

func TestAmfArrayEncode(t *testing.T) {
	testcases, expected := genArrayTestCases()
	codec := NewAmfCodec()
	for i, arr := range testcases {
		// t.Logf("arr[%v]: %v", i, arr)
		result, err := codec.AmfArrayEncode(arr)
		if err != nil {
			t.Errorf("AmfArrayEncode failed: %v", err)
		}
		// codec.LogStringTable("TestAmfArrayEncode")
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfArrayEncode failed: expected:\n%v, got:\n%v", expected[i], result)
		}
	}
}

func TestAmfArrayDecode(t *testing.T) {
	testcases, expected := genArrayTestCases2()
	codec := NewAmfCodec()
	for i, buf := range expected {
		result, _, err := codec.AmfArrayDecode(buf)
		if err != nil {
			t.Errorf("AmfArrayDecode failed: %v", err)
		}
		// codec.LogStringTable("TestAmfArrayDecode")
		if !amfArrayEqual(result, testcases[i]) {
			t.Errorf("AmfArrayDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}
func genArrayTestCases2() ([]*AmfArray, [][]byte) {
	codec := NewAmfCodec()
	str1 := "hello"
	str2 := "world"
	str3 := "hello"
	str4 := "helloworld"
	int1 := uint32(1)
	int2 := uint32(2)
	int3 := uint32(3)
	double1 := 1.0
	double2 := 2.0
	double3 := 3.0
	false1 := false
	false2 := false
	true1 := true
	xmldoc := "<xml></xml>"

	arr1 := EmptyAmfArray()
	arr1.AddDense(str1, AMF_STRING)
	arr1.AddAssoc("key1", int2, AMF_INTEGER)
	arr1.AddAssoc("key2", int3, AMF_INTEGER)
	arr1.AddDense(double1, AMF_DOUBLE)
	arr1.AddDense(double2, AMF_DOUBLE)
	buf1, _ := codec.AmfArrayEncode(arr1)

	arr2 := EmptyAmfArray()
	arr2.AddDense(str2, AMF_STRING)
	arr2.AddAssoc("key1", double2, AMF_DOUBLE)
	arr2.AddAssoc("key2", double2, AMF_DOUBLE)
	arr2.AddDense(int1, AMF_INTEGER)
	arr2.AddDense(int2, AMF_INTEGER)
	buf2, _ := codec.AmfArrayEncode(arr2)

	arr3 := EmptyAmfArray()
	arr3.AddDense(str3, AMF_STRING)
	arr3.AddAssoc("key1", int3, AMF_INTEGER)
	arr3.AddDense(str4, AMF_STRING)
	arr3.AddDense(false1, AMF_FALSE)
	arr3.AddDense(true1, AMF_TRUE)
	arr3.AddAssoc("key2", double3, AMF_DOUBLE)
	buf3, _ := codec.AmfArrayEncode(arr3)

	arr4 := EmptyAmfArray()
	arr4.AddDense(nil, AMF_UNDEFINED)
	arr4.AddDense(nil, AMF_NULL)
	arr4.AddDense(xmldoc, AMF_XML_DOC)
	arr4.AddAssoc("key1", false2, AMF_FALSE)
	buf4, _ := codec.AmfArrayEncode(arr4)

	arr5 := EmptyAmfArray()
	arr5.AddDense(str3, AMF_STRING)
	arr5.AddAssoc("key1", int3, AMF_INTEGER)
	arr5.AddDense(str4, AMF_STRING)
	arr5.AddDense(false1, AMF_FALSE)
	arr5.AddDense(true1, AMF_TRUE)
	arr5.AddAssoc("key2", double3, AMF_DOUBLE)
	buf5, _ := codec.AmfArrayEncode(arr5)

	arrs := []*AmfArray{arr1, arr2, arr3, arr4, arr5}
	bufs := [][]byte{buf1, buf2, buf3, buf4, buf5}
	return arrs, bufs
}

func genArrayTestCases() ([]*AmfArray, [][]byte) {
	codec := NewAmfCodec()
	str1 := "hello"
	str2 := "world"
	str3 := "hello"
	str4 := "helloworld"
	int1 := uint32(1)
	int1buf, _ := AmfIntEncode(int1)
	int2 := uint32(2)
	int2buf, _ := AmfIntEncode(int2)
	// fmt.Println("int2buf:", int2buf)
	int3 := uint32(3)
	int3buf, _ := AmfIntEncode(int3)
	// fmt.Println("int3buf:", int3buf)
	double1 := 1.0
	double1buf, _ := AmfDoubleEncode(double1)
	double2 := 2.0
	double2buf, _ := AmfDoubleEncode(double2)
	// fmt.Println("double2buf:", double2buf)
	double3 := 3.0
	double3buf, _ := AmfDoubleEncode(double3)
	// fmt.Println("double3buf:", double3buf)
	false1 := false
	false1buf := AmfBool(false1)
	false2 := false
	false2buf := AmfBool(false2)
	true1 := true
	true1buf := AmfBool(true1)
	xmldoc := "<xml></xml>"
	undi := AmfUndefined()
	null := AmfNull()

	key1 := "key1"
	key2 := "key2"

	arr1 := EmptyAmfArray()
	arr1.AddDense(str1, AMF_STRING)
	arr1.AddAssoc("key1", int2, AMF_INTEGER)
	arr1.AddAssoc("key2", int3, AMF_INTEGER)
	arr1.AddDense(double1, AMF_DOUBLE)
	arr1.AddDense(double2, AMF_DOUBLE)
	h1_1, _ := codec.AmfStringEncodePayload(key1)
	h1_2, _ := codec.AmfStringEncodePayload(key2)
	h1_3, _ := codec.AmfStringEncode(str1)
	u29a1, _ := AmfIntEncodePayload(3 << 1 | 1)
	assoc1_1 := append(h1_1, int2buf...)
	assoc1_2 := append(h1_2, int3buf...)
	buf1 := append([]byte{AMF_ARRAY}, u29a1...)
	buf1 = append(buf1, assoc1_1...)
	buf1 = append(buf1, assoc1_2...)
	buf1 = append(buf1, 0x01)
	buf1 = append(buf1, h1_3...)
	buf1 = append(buf1, double1buf...)
	buf1 = append(buf1, double2buf...)

	arr2 := EmptyAmfArray()
	arr2.AddDense(str2, AMF_STRING)
	arr2.AddAssoc("key1", double2, AMF_DOUBLE)
	arr2.AddAssoc("key2", double2, AMF_DOUBLE)
	arr2.AddDense(int1, AMF_INTEGER)
	arr2.AddDense(int2, AMF_INTEGER)
	u29a2, _ := AmfIntEncodePayload(3 << 1 | 1)
	r2_1, _ := codec.AmfStringEncodePayload(key1)
	r2_2, _ := codec.AmfStringEncodePayload(key2)
	r2_3, _ := codec.AmfStringEncode(str2)
	assoc2_1 := append(r2_1, double2buf...)
	assoc2_2 := append(r2_2, double2buf...)
	buf2 := append([]byte{AMF_ARRAY}, u29a2...)
	buf2 = append(buf2, assoc2_1...)
	buf2 = append(buf2, assoc2_2...)
	buf2 = append(buf2, 0x01)
	buf2 = append(buf2, r2_3...)
	buf2 = append(buf2, int1buf...)
	buf2 = append(buf2, int2buf...)

	arr3 := EmptyAmfArray()
	arr3.AddDense(str3, AMF_STRING)
	arr3.AddAssoc("key1", int3, AMF_INTEGER)
	arr3.AddDense(str4, AMF_STRING)
	arr3.AddDense(false1, AMF_FALSE)
	arr3.AddDense(true1, AMF_TRUE)
	arr3.AddAssoc("key2", double3, AMF_DOUBLE)
	u29a3, _ := AmfIntEncodePayload(4 << 1 | 1)
	r3_1, _ := codec.AmfStringEncodePayload(key1)
	r3_2, _ := codec.AmfStringEncodePayload(key2)
	h3_3, _ := codec.AmfStringEncode(str3)
	h3_4, _ := codec.AmfStringEncode(str4)
	assoc3_1 := append(r3_1, int3buf...)
	assoc3_2 := append(r3_2, double3buf...)
	buf3 := append([]byte{AMF_ARRAY}, u29a3...)
	buf3 = append(buf3, assoc3_1...)
	buf3 = append(buf3, assoc3_2...)
	buf3 = append(buf3, 0x01)
	buf3 = append(buf3, h3_3...)
	buf3 = append(buf3, h3_4...)
	buf3 = append(buf3, false1buf...)
	buf3 = append(buf3, true1buf...)

	arr4 := EmptyAmfArray()
	arr4.AddDense(nil, AMF_UNDEFINED)
	arr4.AddDense(nil, AMF_NULL)
	arr4.AddDense(xmldoc, AMF_XML_DOC)
	arr4.AddAssoc("key1", false2, AMF_FALSE)
	u29a4, _ := AmfIntEncodePayload(3 << 1 | 1)
	r4_1, _ := codec.AmfStringEncodePayload(key1)
	h4_2, _ := codec.AmfXmlDocEncode(xmldoc)
	assoc4_1 := append(r4_1, false2buf...)
	buf4 := append([]byte{AMF_ARRAY}, u29a4...)
	buf4 = append(buf4, assoc4_1...)
	buf4 = append(buf4, 0x01)
	buf4 = append(buf4, undi...)
	buf4 = append(buf4, null...)
	buf4 = append(buf4, h4_2...)

	arr5 := EmptyAmfArray()
	arr5.AddDense(str3, AMF_STRING)
	arr5.AddAssoc("key1", int3, AMF_INTEGER)
	arr5.AddDense(str4, AMF_STRING)
	arr5.AddDense(false1, AMF_FALSE)
	arr5.AddDense(true1, AMF_TRUE)
	arr5.AddAssoc("key2", double3, AMF_DOUBLE)
	r5_1, _ := codec.AmfStringEncodePayload(key1)
	r5_2, _ := codec.AmfStringEncodePayload(key2)
	r5_3, _ := codec.AmfStringEncode(str3)
	r5_4, _ := codec.AmfStringEncode(str4)
	assoc5_1 := append(r5_1, int3buf...)
	assoc5_2 := append(r5_2, double3buf...)
	buf5 := append([]byte{AMF_ARRAY}, u29a3...)
	buf5 = append(buf5, assoc5_1...)
	buf5 = append(buf5, assoc5_2...)
	buf5 = append(buf5, 0x01)
	buf5 = append(buf5, r5_3...)
	buf5 = append(buf5, r5_4...)
	buf5 = append(buf5, false1buf...)
	buf5 = append(buf5, true1buf...)

	arrs := []*AmfArray{arr1, arr2, arr3, arr4, arr5}
	bufs := [][]byte{buf1, buf2, buf3, buf4, buf5}
	return arrs, bufs
}

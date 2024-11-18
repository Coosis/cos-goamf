package v3

import (
	"fmt"
	"math"
	"testing"
	"encoding/binary"
)

func int32SliceEqual(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range(len(a)) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func uint32SliceEqual(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range(len(a)) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func float64SliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range(len(a)) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func objEqual(a, b *AmfVectorObj) bool {
	if a.TypeName != b.TypeName {
		return false
	}
	if len(a.Data) != len(b.Data) {
		return false
	}
	for i := range(len(a.Data)) {
		if a.Data[i] != b.Data[i] {
			fmt.Printf("i: %v\n", i)
			fmt.Printf("%T:%T", a.Data[i], b.Data[i])
			fmt.Printf("%v:%v", a.Data[i], b.Data[i])
			return false
		}
	}
	if a.FixedLen != b.FixedLen {
		return false
	}
	return true
}

func TestVecIntEncodeDecode(t *testing.T) {
	cs1 := &AmfVectorInt{
		FixedLen: false,
		Data: []int32{0, 1, 2, -3, 4, 5},
	}
	cs2 := &AmfVectorInt{
		FixedLen: false,
		Data: []int32{},
	}
	cs3 := &AmfVectorInt{
		FixedLen: true,
		Data: []int32{1},
	}
	cs4 := &AmfVectorInt{
		FixedLen: true,
		Data: []int32{1, 5, 9, -13},
	}

	testcases := []*AmfVectorInt{
		cs1,
		cs2,
		cs3,
		cs4,
		cs1,
	}

	u29v1, _ := AmfIntEncodePayload(6 << 1 | 1)
	buf1 := []byte{AMF_VECTOR_INT}
	buf1 = append(buf1, u29v1...)
	buf1 = append(buf1, 0x00)
	nums := make([]byte, 6 * 4)
	for i, num := range testcases[0].Data {
		binary.BigEndian.PutUint32(nums[i * 4:], uint32(num))
	}
	buf1 = append(buf1, nums...)

	u29v2, _ := AmfIntEncodePayload(0 << 1 | 1)
	buf2 := []byte{AMF_VECTOR_INT}
	buf2 = append(buf2, u29v2...)
	buf2 = append(buf2, 0x00)

	u29v3, _ := AmfIntEncodePayload(1 << 1 | 1)
	buf3 := []byte{AMF_VECTOR_INT}
	buf3 = append(buf3, u29v3...)
	buf3 = append(buf3, 0x01)
	nums2 := make([]byte, 1 * 4)
	binary.BigEndian.PutUint32(nums2, 1)
	buf3 = append(buf3, nums2...)

	u29v4, _ := AmfIntEncodePayload(4 << 1 | 1)
	buf4 := []byte{AMF_VECTOR_INT}
	buf4 = append(buf4, u29v4...)
	buf4 = append(buf4, 0x01)
	nums3 := make([]byte, 4 * 4)
	for i, num := range testcases[3].Data {
		binary.BigEndian.PutUint32(nums3[i * 4:], uint32(num))
	}
	buf4 = append(buf4, nums3...)

	u29v5, _ := AmfIntEncodePayload(0 << 1)
	buf5 := []byte{AMF_VECTOR_INT}
	buf5 = append(buf5, u29v5...)

	expected := [][]byte{
		buf1,
		buf2,
		buf3,
		buf4,
		buf5,
	}

	codec := NewAmfCodec()
	t.Run("AmfVectorIntEncode", func(t *testing.T) {
		for i, data := range testcases[:2] {
			result, err := codec.AmfVectorIntEncode(data)
			if err != nil {
				t.Errorf("AmfVectorUintEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i]) {
				t.Errorf("AmfVectorUintEncode failed: expected\n%v,\ngot\n%v", expected[i], result)
			}
		}
		for i, data := range testcases[2:] {
			// t.Logf("i: %v", i)
			result, err := codec.AmfVectorIntEncode(data)
			if err != nil {
				t.Errorf("AmfVectorUintEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i+2]) {
				t.Errorf("AmfVectorUintEncode failed: expected\n%v,\ngot\n%v", expected[i+2], result)
			}
		}
	})

	t.Run("AmfVectorIntDecode", func(t *testing.T) {
		for i, data := range expected {
			result, _, err := codec.AmfVectorIntDecode(data)
			if err != nil {
				t.Errorf("AmfVectorIntDecode failed: %v", err)
			}
			if !int32SliceEqual(result.Data, testcases[i].Data) {
				t.Errorf("AmfVectorIntDecode failed: expected\n%v,\ngot\n%v", testcases[i], result)
			}
		}
	})
}

func TestVecUintEncodeDecode(t *testing.T) {
	cs1 := &AmfVectorUint{
		FixedLen: false,
		Data: []uint32{0, 1, 2, 3, 4, 5},
	}
	cs2 := &AmfVectorUint{
		FixedLen: false,
		Data: []uint32{},
	}
	cs3 := &AmfVectorUint{
		FixedLen: true,
		Data: []uint32{1},
	}
	cs4 := &AmfVectorUint{
		FixedLen: true,
		Data: []uint32{1, 5, 9, 13},
	}
	testcases := []*AmfVectorUint{
		cs1,
		cs2,
		cs3,
		cs4,
		cs1,
	}

	u29v1, _ := AmfIntEncodePayload(6 << 1 | 1)
	buf1 := []byte{AMF_VECTOR_UINT}
	buf1 = append(buf1, u29v1...)
	buf1 = append(buf1, 0x00)
	nums := make([]byte, 6 * 4)
	for i, num := range testcases[0].Data {
		binary.BigEndian.PutUint32(nums[i * 4:], num)
	}
	buf1 = append(buf1, nums...)

	u29v2, _ := AmfIntEncodePayload(0 << 1 | 1)
	buf2 := []byte{AMF_VECTOR_UINT}
	buf2 = append(buf2, u29v2...)
	buf2 = append(buf2, 0x00)

	u29v3, _ := AmfIntEncodePayload(1 << 1 | 1)
	buf3 := []byte{AMF_VECTOR_UINT}
	buf3 = append(buf3, u29v3...)
	buf3 = append(buf3, 0x01)
	nums2 := make([]byte, 1 * 4)
	binary.BigEndian.PutUint32(nums2, 1)
	buf3 = append(buf3, nums2...)

	u29v4, _ := AmfIntEncodePayload(4 << 1 | 1)
	buf4 := []byte{AMF_VECTOR_UINT}
	buf4 = append(buf4, u29v4...)
	buf4 = append(buf4, 0x01)
	nums3 := make([]byte, 4 * 4)
	for i, num := range testcases[3].Data {
		binary.BigEndian.PutUint32(nums3[i * 4:], num)
	}
	buf4 = append(buf4, nums3...)

	u29v5, _ := AmfIntEncodePayload(0 << 1)
	buf5 := []byte{AMF_VECTOR_UINT}
	buf5 = append(buf5, u29v5...)

	expected := [][]byte{
		buf1,
		buf2,
		buf3,
		buf4,
		buf5,
	}

	codec := NewAmfCodec()
	t.Run("AmfVectorUintEncode", func(t *testing.T) {
		for i, data := range testcases[:2] {
			result, err := codec.AmfVectorUintEncode(data)
			if err != nil {
				t.Errorf("AmfVectorUintEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i]) {
				t.Errorf("AmfVectorUintEncode failed: expected\n%v,\ngot\n%v", expected[i], result)
			}
		}
		for i, data := range testcases[2:] {
			// t.Logf("i: %v", i)
			result, err := codec.AmfVectorUintEncode(data)
			if err != nil {
				t.Errorf("AmfVectorUintEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i+2]) {
				t.Errorf("AmfVectorUintEncode failed: expected\n%v,\ngot\n%v", expected[i+2], result)
			}
		}
	})

	t.Run("AmfVectorUintDecode", func(t *testing.T) {
		for i, data := range expected {
			result, _, err := codec.AmfVectorUintDecode(data)
			if err != nil {
				t.Errorf("AmfVectorUintDecode failed: %v", err)
			}
			if !uint32SliceEqual(result.Data, testcases[i].Data) {
				t.Errorf("AmfVectorUintDecode failed: expected\n%v,\ngot\n%v", testcases[i], result)
			}
		}
	})
}

func TestVecDoubleEncodeDecode(t *testing.T) {
	cs1 := &AmfVectorDouble{
		FixedLen: false,
		Data: []float64{0, 1.2, 2.4, 3.5, 4.11, 5.6321},
	}
	cs2 := &AmfVectorDouble{
		FixedLen: false,
		Data: []float64{},
	}
	cs3 := &AmfVectorDouble{
		FixedLen: true,
		Data: []float64{1},
	}
	cs4 := &AmfVectorDouble{
		FixedLen: true,
		Data: []float64{1.9, 5.2, 9.2342434, 13.0},
	}
	testcases := []*AmfVectorDouble{
		cs1,
		cs2,
		cs3,
		cs4,
		cs1,
	}

	u29v1, _ := AmfIntEncodePayload(6 << 1 | 1)
	buf1 := []byte{AMF_VECTOR_DOUBLE}
	buf1 = append(buf1, u29v1...)
	buf1 = append(buf1, 0x00)
	nums := make([]byte, 6 * 8)
	for i, num := range testcases[0].Data {
		binary.BigEndian.PutUint64(nums[i * 8:], math.Float64bits(num))
	}
	buf1 = append(buf1, nums...)

	u29v2, _ := AmfIntEncodePayload(0 << 1 | 1)
	buf2 := []byte{AMF_VECTOR_DOUBLE}
	buf2 = append(buf2, u29v2...)
	buf2 = append(buf2, 0x00)

	u29v3, _ := AmfIntEncodePayload(1 << 1 | 1)
	buf3 := []byte{AMF_VECTOR_DOUBLE}
	buf3 = append(buf3, u29v3...)
	buf3 = append(buf3, 0x01)
	nums2 := make([]byte, 1 * 8)
	binary.BigEndian.PutUint64(nums2, math.Float64bits(1))
	buf3 = append(buf3, nums2...)

	u29v4, _ := AmfIntEncodePayload(4 << 1 | 1)
	buf4 := []byte{AMF_VECTOR_DOUBLE}
	buf4 = append(buf4, u29v4...)
	buf4 = append(buf4, 0x01)
	nums3 := make([]byte, 4 * 8)
	for i, num := range testcases[3].Data {
		binary.BigEndian.PutUint64(nums3[i * 8:], math.Float64bits(num))
	}
	buf4 = append(buf4, nums3...)

	u29v5, _ := AmfIntEncodePayload(0 << 1)
	buf5 := []byte{AMF_VECTOR_DOUBLE}
	buf5 = append(buf5, u29v5...)

	expected := [][]byte{
		buf1,
		buf2,
		buf3,
		buf4,
		buf5,
	}

	codec := NewAmfCodec()
	t.Run("AmfVectorDoubleEncode", func(t *testing.T) {
		for i, data := range testcases[:2] {
			result, err := codec.AmfVectorDoubleEncode(data)
			if err != nil {
				t.Errorf("AmfVectorDoubleEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i]) {
				t.Errorf("AmfVectorDoubleEncode failed: expected\n%v,\ngot\n%v", expected[i], result)
			}
		}
		for i, data := range testcases[2:] {
			// t.Logf("i: %v", i)
			result, err := codec.AmfVectorDoubleEncode(data)
			if err != nil {
				t.Errorf("AmfVectorDoubleEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i+2]) {
				t.Errorf("AmfVectorDoubleEncode failed: expected\n%v,\ngot\n%v", expected[i+2], result)
			}
		}
	})

	t.Run("AmfVectorDoubleDecode", func(t *testing.T) {
		for i, data := range expected {
			// fmt.Println(i)
			result, _, err := codec.AmfVectorDoubleDecode(data)
			if err != nil {
				t.Errorf("AmfVectorDoubleDecode failed: %v", err)
			}
			if !float64SliceEqual(result.Data, testcases[i].Data) {
				t.Errorf("AmfVectorDoubleDecode failed: expected\n%v,\ngot\n%v", testcases[i], result)
			}
		}
	})
}

func TestVecObjEncodeDecode(t *testing.T) {
	obj1 := &AmfVectorObj{
		TypeName: "int",
		FixedLen: false,
		Data: []interface{}{uint32(1), uint32(2), uint32(3)},
	}
	obj2 := &AmfVectorObj{
		TypeName: "int",
		FixedLen: false,
		Data: []interface{}{},
	}
	obj3 := &AmfVectorObj{
		TypeName: "bool",
		FixedLen: true,
		Data: []interface{}{true},
	}
	obj4 := &AmfVectorObj{
		TypeName: "xml",
		FixedLen: true,
		Data: []interface{}{
			AmfXml("<xml><img></img></xml>"),
			AmfXml("<xml></xml>"),
			AmfXml("<img></img>"),
			AmfXml("<xml><img></img></xml>"),
		},
	}

	codec := NewAmfCodec()
	u29v1, _ := AmfIntEncodePayload(3 << 1 | 1)
	buf1 := []byte{AMF_VECTOR_OBJECT}
	buf1 = append(buf1, u29v1...)
	buf1 = append(buf1, 0x00)
	name1, _ := codec.AmfStringEncodePayload("int")
	buf1 = append(buf1, name1...)
	int1, _ := AmfIntEncode(1)
	int2, _ := AmfIntEncode(2)
	int3, _ := AmfIntEncode(3)
	buf1 = append(buf1, int1...)
	buf1 = append(buf1, int2...)
	buf1 = append(buf1, int3...)
	codec.Append(&obj1, COMPLEX_TABLE)

	u29v2, _ := AmfIntEncodePayload(0 << 1 | 1)
	buf2 := []byte{AMF_VECTOR_OBJECT}
	buf2 = append(buf2, u29v2...)
	buf2 = append(buf2, 0x00)
	name2, _ := codec.AmfStringEncodePayload("int")
	buf2 = append(buf2, name2...)
	codec.Append(&obj2, COMPLEX_TABLE)

	u29v3, _ := AmfIntEncodePayload(1 << 1 | 1)
	buf3 := []byte{AMF_VECTOR_OBJECT}
	buf3 = append(buf3, u29v3...)
	buf3 = append(buf3, 0x01)
	name3, _ := codec.AmfStringEncodePayload("bool")
	buf3 = append(buf3, name3...)
	buf3 = append(buf3, AMF_TRUE)
	codec.Append(&obj3, COMPLEX_TABLE)

	u29v4, _ := AmfIntEncodePayload(4 << 1 | 1)
	buf4 := []byte{AMF_VECTOR_OBJECT}
	buf4 = append(buf4, u29v4...)
	buf4 = append(buf4, 0x01)
	name4, _ := codec.AmfStringEncodePayload("xml")
	buf4 = append(buf4, name4...)
	xml1, _ := codec.AmfXmlEncode("<xml><img></img></xml>")
	xml2, _ := codec.AmfXmlEncode("<xml></xml>")
	xml3, _ := codec.AmfXmlEncode("<img></img>")
	xml4, _ := codec.AmfXmlEncode("<xml><img></img></xml>")
	buf4 = append(buf4, xml1...)
	buf4 = append(buf4, xml2...)
	buf4 = append(buf4, xml3...)
	buf4 = append(buf4, xml4...)
	codec.Append(&obj4, COMPLEX_TABLE)

	u29v5, _ := AmfIntEncodePayload(0 << 1)
	buf5 := []byte{AMF_VECTOR_OBJECT}
	buf5 = append(buf5, u29v5...)

	expected := [][]byte{
		buf1,
		buf2,
		buf3,
		buf4,
		buf5,
	}

	t.Run("AmfVectorObjEncode", func(t *testing.T) {
		codec := NewAmfCodec()
		r1, err := codec.AmfVectorObjEncode(obj1)
		if err != nil {
			t.Errorf("AmfVectorObjEncode failed: %v", err)
		}
		if !bytesEqual(r1, expected[0]) {
			t.Errorf("AmfVectorObjEncode failed: expected\n%v,\ngot\n%v", expected[0], r1)
		}

		r2, err := codec.AmfVectorObjEncode(obj2)
		if err != nil {
			t.Errorf("AmfVectorObjEncode failed: %v", err)
		}
		if !bytesEqual(r2, expected[1]) {
			t.Errorf("AmfVectorObjEncode failed: expected\n%v,\ngot\n%v", expected[1], r2)
		}

		r3, err := codec.AmfVectorObjEncode(obj3)
		if err != nil {
			t.Errorf("AmfVectorObjEncode failed: %v", err)
		}
		if !bytesEqual(r3, expected[2]) {
			t.Errorf("AmfVectorObjEncode failed: expected\n%v,\ngot\n%v", expected[2], r3)
		}

		r4, err := codec.AmfVectorObjEncode(obj4)
		if err != nil {
			t.Errorf("AmfVectorObjEncode failed: %v", err)
		}
		if !bytesEqual(r4, expected[3]) {
			t.Errorf("AmfVectorObjEncode failed: expected\n%v,\ngot\n%v", expected[3], r4)
		}

		r5, err := codec.AmfVectorObjEncode(obj1)
		if err != nil {
			t.Errorf("AmfVectorObjEncode failed: %v", err)
		}
		if !bytesEqual(r5, expected[4]) {
			t.Errorf("AmfVectorObjEncode failed: expected\n%v,\ngot\n%v", expected[4], r5)
		}
	})

	t.Run("AmfVectorObjDecode", func(t *testing.T) {
		codec := NewAmfCodec()
		r1, _, err := codec.AmfVectorObjDecode(expected[0])
		if err != nil {
			t.Errorf("AmfVectorObjDecode failed: %v", err)
		}
		if !objEqual(r1, obj1) {
			t.Errorf("AmfVectorObjDecode failed: expected\n%v,\ngot\n%v", obj1, r1)
		}

		r2, _, err := codec.AmfVectorObjDecode(expected[1])
		if err != nil {
			t.Errorf("AmfVectorObjDecode failed: %v", err)
		}
		if !objEqual(r2, obj2) {
			t.Errorf("AmfVectorObjDecode failed: expected\n%v,\ngot\n%v", obj2, r2)
		}

		r3, _, err := codec.AmfVectorObjDecode(expected[2])
		if err != nil {
			t.Errorf("AmfVectorObjDecode failed: %v", err)
		}
		if !objEqual(r3, obj3) {
			t.Errorf("AmfVectorObjDecode failed: expected\n%v,\ngot\n%v", obj3, r3)
		}

		r4, _, err := codec.AmfVectorObjDecode(expected[3])
		if err != nil {
			t.Errorf("AmfVectorObjDecode failed: %v", err)
		}
		if !objEqual(r4, obj4) {
			t.Errorf("AmfVectorObjDecode failed: expected\n%v,\ngot\n%v", obj4, r4)
		}

		r5, _, err := codec.AmfVectorObjDecode(expected[4])
		if err != nil {
			t.Errorf("AmfVectorObjDecode failed: %v", err)
		}
		if !objEqual(r5, obj1) {
			t.Errorf("AmfVectorObjDecode failed: expected\n%v,\ngot\n%v", obj1, r5)
		}
	})
}

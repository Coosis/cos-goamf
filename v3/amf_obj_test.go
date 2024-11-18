package v3

import (
	"fmt"
	"testing"
)

func amfObjEqual(a, b *AmfObj) bool {
	if a.ClassName != b.ClassName {
		return false
	}
	if len(a.Member) != len(b.Member) {
		return false
	}
	if len(a.DynMembers) != len(b.DynMembers) {
		return false
	}
	for i := range a.Member {
		av := a.Member[i]
		bv := b.Member[i]
		if av.Key != bv.Key {
			return false
		}
		if !valuesEqual(av.Value, bv.Value) {
			fmt.Printf("%T:%T", av.Value, bv.Value)
			return false
		}
	}
	for i := range a.DynMembers {
		av := a.DynMembers[i]
		bv := b.DynMembers[i]
		if av.Key != bv.Key {
			return false
		}
		if !valuesEqual(av.Value, bv.Value) {
			return false
		}
	}
	if !bytesEqual(a.ExtTraits, b.ExtTraits) {
		return false
	}
	return true
}

func genObjTestcases() ([]*AmfObj, [][]byte) {
	name1 := "name1"
	name2 := "name2"
	name3 := "name3"

	int2 := uint32(2)
	int3 := uint32(3)
	double1 := 3.0
	false1 := false
	xmldoc := AmfXmlDoc("<xml></xml>")

	// traits
	obj1 := EmptyAmfObj()
	obj1.ClassName = name1
	obj1.AppendMember(AmfObjMember{Key: "key1", Value: int2})
	obj1.AppendMember(AmfObjMember{Key: "key2", Value: int3})
	obj1.AppendDynMember(AmfObjMember{Key: "key1", Value: double1})

	// traits
	obj2 := EmptyAmfObj()
	obj2.ClassName = name2
	obj2.AppendMember(AmfObjMember{Key: "key1", Value: int3})
	obj2.AppendMember(AmfObjMember{Key: "key2", Value: xmldoc})
	obj2.AppendDynMember(AmfObjMember{Key: "key1", Value: false1})

	// traits-ext
	obj3 := EmptyAmfObj()
	obj3.ClassName = name3
	obj3.ExtTraits = []byte{0x01, 0x02, 0x03}

	objs := []*AmfObj{obj1, obj2, obj1, obj3}

	codec := NewAmfCodec()

	h1, _ := AmfIntEncodePayload((2 << 4) | 0x0b)
	buf1 := []byte{AMF_OBJECT}
	buf1 = append(buf1, h1...)
	// fmt.Println(buf1)
	classname1, _ := codec.AmfStringEncodePayload(name1)
	buf1 = append(buf1, classname1...)
	// fmt.Println(buf1)
	traitname11, _ := codec.AmfStringEncodePayload("key1")
	traitname12, _ := codec.AmfStringEncodePayload("key2")
	buf1 = append(buf1, traitname11...)
	// fmt.Println(buf1)
	buf1 = append(buf1, traitname12...)
	// fmt.Println(buf1)
	int11buf, _ := AmfIntEncode(int2)
	int12buf, _ := AmfIntEncode(int3)
	buf1 = append(buf1, int11buf...)
	// fmt.Println(buf1)
	buf1 = append(buf1, int12buf...)
	// fmt.Println(buf1)
	traitname13, _ := codec.AmfStringEncodePayload("key1")
	double1buf, _ := AmfDoubleEncode(double1)
	buf1 = append(buf1, traitname13...)
	// fmt.Println(buf1)
	buf1 = append(buf1, double1buf...)
	// fmt.Println(buf1)
	buf1 = append(buf1, 0x01)
	// fmt.Println(buf1)

	h2, _ := AmfIntEncodePayload((2 << 4) | 0x0b)
	buf2 := []byte{AMF_OBJECT}
	buf2 = append(buf2, h2...)
	// fmt.Println(buf2)
	classname2, _ := codec.AmfStringEncodePayload(name2)
	buf2 = append(buf2, classname2...)
	// fmt.Println(buf2)
	traitname21, _ := codec.AmfStringEncodePayload("key1")
	traitname22, _ := codec.AmfStringEncodePayload("key2")
	buf2 = append(buf2, traitname21...)
	buf2 = append(buf2, traitname22...)
	// fmt.Println(buf2)
	int21buf, _ := AmfIntEncode(int3)
	buf2 = append(buf2, int21buf...)
	// fmt.Println(buf2)
	xmldocbuf, _ := codec.AmfXmlDocEncode(xmldoc)
	buf2 = append(buf2, xmldocbuf...)
	// fmt.Println(buf2)
	traitname23, _ := codec.AmfStringEncodePayload("key1")
	false1buf := []byte{AMF_FALSE}
	buf2 = append(buf2, traitname23...)
	// fmt.Println(buf2)
	buf2 = append(buf2, false1buf...)
	// fmt.Println(buf2)
	buf2 = append(buf2, 0x01)
	// fmt.Println(buf2)

	buf1alt := []byte{AMF_OBJECT}
	ref1, _ := AmfIntEncodePayload(0 << 1)
	buf1alt = append(buf1alt, ref1...)

	h3, _ := AmfIntEncodePayload(0x07)
	buf3 := []byte{AMF_OBJECT}
	buf3 = append(buf3, h3...)
	classname3, _ := codec.AmfStringEncodePayload(name3)
	buf3 = append(buf3, classname3...)
	buf3 = append(buf3, obj3.ExtTraits...)

	return objs, [][]byte{buf1, buf2, buf1alt, buf3}
}

func TestAmfObjEncode(t *testing.T) {
	codec := NewAmfCodec()
	testcases, expected := genObjTestcases()

	t.Run("AmfObjEncode", func(t *testing.T) {
		for i, value := range testcases {
			// t.Logf("i: %v", i)
			result, err := codec.AmfObjEncode(value)
			if err != nil {
				t.Errorf("AmfObjEncode failed: %v", err)
			}
			if !bytesEqual(result, expected[i]) {
				t.Errorf("AmfObjEncode failed: expected:\n%v,\ngot:\n%v", expected[i], result)
			}
		}
	})
}

func TestAmfObjDecode(t *testing.T) {
	codec := NewAmfCodec()
	expected, testcases := genObjTestcases()
	for i, buf := range testcases {
		// t.Logf("i: %v", i)
		result, _, err := codec.AmfObjDecode(buf)
		if err != nil {
			t.Errorf("AmfObjDecode failed: %v", err)
		}
		// if result != expected[i] {
		if !amfObjEqual(expected[i], result) {
			t.Errorf("AmfObjDecode failed: expected\n%v,\ngot\n%v", expected[i], result)
		}
	}
}

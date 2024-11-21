package v0

import (
	"testing"
	"reflect"
)

func objEqual(obj1, obj2 *AmfObj) bool {
	if obj1.Name != obj2.Name {
		return false
	}
	if len(obj1.Props) != len(obj2.Props) {
		return false
	}
	for i := 0; i < len(obj1.Props); i++ {
		if !reflect.DeepEqual(obj1.Props[i], obj2.Props[i]) {
			return false
		}
	}
	return true
}

func genObjTestcases() []*AmfObj {
	obj1 := NewAmfObj()
	obj1.AddProp("key1", "value1")
	obj1.AddProp("key2", float64(1))

	obj2 := NewAmfObj()
	obj2.Name = "obj2"
	obj2.AddProp("key1", true)
	obj2.AddProp("key2", AmfDate(1.1))

	obj3 := NewAmfObj()
	obj3.Name = "obj3"
	obj3.AddProp("key1", AmfXmldoc("<xml></xml>"))
	obj3.AddProp("key2", NewAmfObj())

	return []*AmfObj{obj1, obj2, obj3, obj1, obj2}
}

func TestAmfObjEncodeDecode(t *testing.T) {
	testcases := genObjTestcases()
	c := NewAmfCodec()
	for _, value := range testcases {
		encoded, err := c.AmfObjEncode(value)
		// t.Logf("encoded: %v", encoded)
		if err != nil {
			t.Errorf("AmfObjEncode failed: %v", err)
		}
		decoded, _, err := c.AmfObjDecode(encoded)
		if err != nil {
			t.Errorf("AmfObjDecode failed: %v", err)
		}
		if !objEqual(value, decoded) {
			t.Errorf("AmfObjEncodeDecode failed: expected %v, got %v", value, decoded)
		}
	}
}

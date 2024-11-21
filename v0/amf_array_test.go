package v0

import (
	// "fmt"
	"testing"
	"reflect"
)

func genArrayTestcases() []*AmfArray {
	arr1 := NewAmfArray()
	arr1.Add("hello")
	arr1.Add(AmfDate(1.1))

	arr2 := NewAmfArray()
	arr2.Add(AmfXmldoc("<xml></xml>"))
	arr2.Add(NewAmfObj())

	arr3 := NewAmfArray()
	arr3.Add(true)
	arr3.Add(NewAmfObj())

	return []*AmfArray{arr1, arr2, arr3, arr1}
}

func TestAmfArrayEncodeDecode(t *testing.T) {
	testcases := genArrayTestcases()
	c := NewAmfCodec()
	for _, value := range testcases {
		encoded, err := c.AmfArrayEncode(value)
		if err != nil {
			t.Errorf("AmfArrayEncode failed: %v", err)
		}
		// fmt.Printf("encoded: %v\n", encoded)
		decoded, _, err := c.AmfArrayDecode(encoded)
		if err != nil {
			t.Errorf("AmfArrayDecode failed: %v", err)
		}
		if !reflect.DeepEqual(value, decoded) {
			t.Errorf("AmfArrayEncodeDecode failed: expected %v, got %v", value, decoded)
		}
	}
}

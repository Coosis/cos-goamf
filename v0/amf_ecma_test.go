package v0

import (
	// "fmt"
	"testing"
	"reflect"
)

func ecmaEqual(a, b *AmfECMA) bool {
	if len(*a) != len(*b) {
		// fmt.Print("length not same\n")
		return false
	}

	for key, val := range(*a) {
		if !reflect.DeepEqual(val, (*b)[key]) {
			// fmt.Print("val not same\n")
			// fmt.Printf("%T: %v and %T: %v\n", val, val, (*b)[key], (*b)[key])
			return false
		}
	}

	return true
}

func genECMATestcases() []*AmfECMA {
	// Testcase 1
	ecma1 := NewAmfECMA()
	(*ecma1)["key1"] = "value1"
	(*ecma1)["key2"] = float64(1)

	// Testcase 2
	ecma2 := NewAmfECMA()
	(*ecma2)["key1"] = true
	(*ecma2)["key2"] = AmfDate(1.1)

	// Testcase 3
	ecma3 := NewAmfECMA()
	obj1 := NewAmfObj()
	obj1.AddProp("key1", "value1")
	(*ecma3)["key1"] = obj1
	(*ecma3)["key2"] = AmfXmldoc("<xml></xml>")

	return []*AmfECMA{ecma1, ecma2, ecma3, ecma1}
}

func TestAmfECMAEncodeDecode(t *testing.T) {
	testcases := genECMATestcases()
	c := NewAmfCodec()
	for _, value := range testcases {
		// t.Logf("Testcase %v: %v", i, value)
		encoded, err := c.AmfECMAEncode(value)
		// t.Logf("encoded: %v", encoded)
		if err != nil {
			t.Errorf("AmfECMAEncode failed: %v", err)
		}
		decoded, _, err := c.AmfECMADecode(encoded)
		if err != nil {
			t.Errorf("AmfECMADecode failed: %v", err)
		}
		if !ecmaEqual(value, decoded) {
			t.Errorf("AmfECMAEncodeDecode failed: expected %v, got %v", value, decoded)
		}
	}
}

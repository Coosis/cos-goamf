package main

import (
	"fmt"
    . "github.com/Coosis/cos-goamf/v0"
)

func ExampleAmfNumEncodeDecode() {
	// encode 
	somevar := float64(1)
	encoded, err := AmfNumberEncode(somevar)
	if err != nil {
		return
	}
	fmt.Printf("Encoded %v to: %v\n", somevar, encoded)

	// decode
	decoded, numbytes, err := AmfNumberDecode(encoded)
	if err != nil {
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, decoded)
}

func ExampleAmfStrEncodeDecode() {
	// encode
	somestr := "hello"
	encoded := AmfStringEncode(somestr)
	fmt.Printf("Encoded %v to: %v\n", somestr, encoded)

	// decode
	decoded, numbytes, err := AmfStringDecode(encoded)
	if err != nil {
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, decoded)
	fmt.Printf("Decoded %v from: %v\n", decoded, encoded)
}

func ExampleAmfObjEncodeDecode() {
	// where ref tables are stored
	codec := NewAmfCodec()

	// encode
	obj := NewAmfObj()
	obj.AddProp("key1", "value1")
	obj.AddProp("key2", true)

	encoded, err := codec.AmfObjEncode(obj)
	if err != nil {
		return
	}
	fmt.Printf("Encoded %v to: %v\n", obj, encoded)

	// decode
	decoded, numbytes, err := codec.AmfObjDecode(encoded)
	if err != nil {
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, decoded)

	encoded2, err := codec.AmfObjEncode(obj)
	if err != nil {
		return
	}
	fmt.Printf("Encoded %v to: %v\n", obj, encoded2)

	// decode
	endecoded, numbytes, err := codec.AmfObjDecode(encoded2)
	if err != nil {
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, endecoded)
}

func main() {
	ExampleAmfNumEncodeDecode()
	ExampleAmfStrEncodeDecode()
	ExampleAmfObjEncodeDecode()
}

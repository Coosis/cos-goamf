package main

import (
    "fmt"
    . "github.com/Coosis/cos-goamf/v3"
)

func ExampleAmfIntEncodeDecode() {
	// encode 
	somevar := uint32(1)
	encoded, err := AmfIntEncode(somevar)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Encoded %v to: %v\n", somevar, encoded)

	// decode
	decoded, numbytes, err := AmfIntDecode(encoded)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, decoded)
	fmt.Printf("Decoded %v from: %v\n", decoded, encoded)
}

func ExampleAmfStringEncodeDecode() {
	// where ref tables are stored
    codec := NewAmfCodec()

	// encode
	somestr := "hello"
	somestr2 := "hello"
	encoded, err := codec.AmfStringEncode(somestr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Encoded %v to: %v\n", somestr, encoded)

	encoded2, err := codec.AmfStringEncode(somestr2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Encoded %v to: %v\n", somestr2, encoded2)

	// decode
	decoded, numbytes, err := codec.AmfStringDecode(encoded)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, decoded)
	fmt.Printf("Decoded %v from: %v\n", decoded, encoded)

	decoded2, numbytes, err := codec.AmfStringDecode(encoded2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Read from the first %v bytes: %v\n", numbytes, decoded2)
	fmt.Printf("Decoded %v from: %v\n", decoded2, encoded2)
}

func main() {
    // simple types:
	ExampleAmfIntEncodeDecode()

    // types that requires the ref table:
	ExampleAmfStringEncodeDecode()
}

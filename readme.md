# A go lib for handling amf v0 and v3 that aims to support all the type mentioned in the spec:
## This lib aims to support all the type mentioned in the spec(yes, with the ref table):
### AMF 0:
- number
- boolean
- string
- object
- movieclip (encode and decode not possible)
- null
- undefined (encode not possible)
- reference
- ecma-array
- object end
- strict array
- date
- long string
- unsupported
- recordset
- xml document
- typed object
### AMF 3:
- undefined
- null
- false
- true
- integer
- double
- string
- xml-doc
- date
- array
- object
- xml
- byte array
- vector int
- vector uint
- vector double
- vector object
- dictionary
## Getting Started
```bash
go get github.com/Coosis/cos-goamf
```
To use this in your code:
```go
import "github.com/Coosis/cos-goamf/v0" // for amf v0
import "github.com/Coosis/cos-goamf/v3" // for amf v3
```

## AMF0 Usage
```go
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
```

## AMF3 Usage
```go
package main

import (
    "fmt"
    . "github.com/Coosis/cos-goamf/v3"
    // or use:
    // "github.com/Coosis/cos-goamf/v3"
    // and use v3.xxx to access
)

func main() {
    // simple types:
    // encode 
    encoded, err := AmfIntEncode(uint32(1))
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(encoded)

    // decode
    decoded, numbytes, err := AmfIntDecode(encoded)
    if err != nil {
        fmt.Println(err)
        return
    }
	fmt.Println(numbytes)
    fmt.Println(decoded)

    // types that requires the ref table:
    codec := NewAmfCodec()
    // encode
    str1 := "hello"
    str2 := "hello"
    encoded, err = codec.AmfStringEncode(str1)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(encoded)
	encoded2, err := codec.AmfStringEncode(str2)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(encoded2)

    // decode
	decodedstr, numbytes, err := codec.AmfStringDecode(encoded)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(decodedstr)
	decodedstr2, numbytes, err := codec.AmfStringDecode(encoded2)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(decodedstr2)
}
```

# Additional Notes:
Sometimes you will encounter where you see both xxxencode and xxxencodePayload. As you know, Amf v3 
has a marker byte that indicates the type of the payload. Basically what happens is as follows:
```
Encode: Marker handling + encode payload
Decode: Marker handling + decode payload
```

If you discover any bugs or have any suggestions, feel free to open an issue or a pull request.

# TODO:
[ ] Add more robust tests

[ ] Fix all the "TODO" in the code(feel free to grep and fix, then submit a pr)

[ ] Add documentation

[x] Implement amf v0 as well

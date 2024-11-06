# A go lib for handling amf version 3 that aims to support all the type mentioned in the spec:
## This lib aims to support all the type mentioned in the spec(yes, with the ref table):
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
import "github.com/Coosis/cos-goamf"
```

## Usage
```go
package main

import (
    "fmt"
    "github.com/Coosis/cos-goamf"
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
    decoded, err := amf.Decode(encoded)
    if err != nil {
        fmt.Println(err)
        return
    }
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
    encoded2, err = codec.AmfStringEncode(str2)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(encoded2)

    // decode
    decoded, err = codec.Decode(encoded)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(decoded)
    decoded2, err = codec.Decode(encoded2)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(decoded2)
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
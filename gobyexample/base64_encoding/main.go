// go provides built-in support for base64 encoding/decoding.

package main

// this syntax imports the 'encoding/base64' package with
// the 'b64' name instead of the default 'base64'. it'll
// save us some space below.
import (
	b64 "encoding/base64"
	"fmt"
)

func main() {
	// here's the 'string' we'll encode/decode.
	data := "abc123!?$*&()'-=@~"

	// go supports both standard and URL-compatible
	// base64. here's how to encode using the standard
	// encoder. the encoder requires a '[]byte' so we
	// cast our 'string' to that type.
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	// decoding may return an error, which you can check
	// if you don't already know the input to be
	// well-formed.
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
	fmt.Println()

	// this encodes/decodes using a URL-compatible base64
	// format.
	uEnc := b64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}

package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

func main() {
	fmt.Println("Cgo must always be guarded with build tags.")
	const s = "48656c6c6f20476f7068657221"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", decoded)

	// Creating UUID Version 4
	// panic on error
	u1, _ := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)

	// or error handling
	u2, _ := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u2)

	// Parsing UUID from string input
	u3, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s\n", u3)

	fmt.Printf("hex: %s\n", u3.Bytes())
	fmt.Printf("encode: %s\n", hex.EncodeToString(u3.Bytes()))

	tmpDst := make([]byte, hex.DecodedLen(len(u3.Bytes())))
	tmpn, err := hex.Decode(tmpDst, u3.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("res: %s\n", tmpDst[:tmpn])

	src := []byte("48656c6c6f20476f7068657221")

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("res: %s\n", dst[:n])
}

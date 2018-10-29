// sha1 hashes are
// frequently used to compute short identies for binary
// or text blobs. for example, the git revision control
// system uses sha1s extensively to
// identify versioned files and directories. here's how to
// compute sha1 hashes in go

package main

// go implements several hash functions in various
// 'crypto/*' packages.
import (
	"crypto/sha1"
	"fmt"
)

func main() {
	s := "sha1 this string"

	// the pattern for generating a hash is 'sha1.New()'.
	// 'sha1.Write(bytes)', then 'sha1.Sum([]byte{})'.
	// here we start with a new hash.
	h := sha1.New()

	// 'Write' expects bytes. If you have a string 's',
	// use '[]byte(s)' to coerce it to bytes.
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. the argument to 'Sum' can be used to append
	// to an existing byte clise: it usually isn't needed.
	bs := h.Sum(nil)

	// sha1 values are often printed in hex, for example
	// in git commits. use the '%x' format verb to convert
	// a hash results to a hex string.
	fmt.Println(s)
	fmt.Printf("%x\n", bs)
}

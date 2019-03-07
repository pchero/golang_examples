// defer is used to ensure that a function call is
// performed later in a program's execution, usually for
// purpose of cleanup. 'defer' is often used where e.g.
// 'ensure' and 'finally' would be used in other languages.

package main

import "fmt"
import "os"

// suppose we wanted to create a file, write to it,
// and then close when we're done. here's how we could
// do that with 'defer'.
func main() {
	// immediately after getting a file object with
	// 'createFile', we defer the closing of that file
	// with 'closeFile'. this will be executed at the end
	// of the enclosing function ('main'), after
	// 'writeFile' has finished.
	f := createFile("/tmp/defer.txt")
	defer closeFile(f)
	writeFile(f)
}

func createFile(p string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

func writeFile(f *os.File) {
	fmt.Println("Writing")
	fmt.Fprintln(f, "data")
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	f.Close()
}

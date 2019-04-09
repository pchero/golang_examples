// https://gobyexample.com/spawning-processes

// Sometimes our Go programs need to spawn other, non-Go
// processes. For example, the syntax highlighting on this
// site is implemented
// by spawning a 'pygmentize'
// process from a Go prgoram. Let's look at a few examples
// of spawning processes from Go.

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	// We'll start with a simple command that takes no
	// arguments or input and just prints something to
	// stdout. The 'exec.Command' helper creates an object
	// to represent this external process.
	dateCmd := exec.Command("date")

	// '.Output' is another helper that handles the common
	// case of running a command, waiting for it to finishe,
	// and collecting its output. If there were no errors,
	// 'dateOut' will hold bytes with the date info.
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))

	// Next we'll look at a slightly more involved case
	// where we pipe date to the external process on its
	// 'stdin' and collect theresult from its 'stdout'.
	grepCmd := exec.Command("grep", "hello")

	// Here we explicitly grab input/output pipes, start
	// the process, write some input to it, read the
	// resulting output, and finally wait for the process
	// to exit.
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	grepCmd.Start()

	grepIn.Write([]byte("hello grep\ngoobye grep"))
	grepIn.Close()

	grepBytes, _ := ioutil.ReadAll(grepOut)
	grepCmd.Wait()

	// We ommited error checks in the above example, but
	// you could use the usaul 'if err!= nil' pattern for
	// all of them. We also only collect the 'StdoutPipe'
	// results, but you could collect the 'StderrPipe' in
	// exactly the same way.
	fmt.Println("> grep hello")
	fmt.Println(string(grepBytes))

	// Note that when spawning commands we need to
	// provided an explicitly delineated command and
	// argument array, vs. being able to just pass in one
	// command-line string. If you want to spawn a full
	// command with a string, you can use `bash` 's `-c`
	// option:
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}

// $ go run cmd/gobyexample_2/spawning_processes/main.go
// > date
// Tue Apr  9 09:13:45 CEST 2019

// > grep hello
// hello grep

// > ls -a -l -h
// total 40
// drwxr-xr-x   8 sungtaekim  staff   256B Mar 14 09:17 .
// drwxr-xr-x   4 sungtaekim  staff   128B Mar 14 09:17 ..
// -rw-r--r--@  1 sungtaekim  staff   6.0K Mar 14 09:17 .DS_Store
// drwxr-xr-x  15 sungtaekim  staff   480B Apr  9 09:13 .git
// -rw-r--r--   1 sungtaekim  staff    17B Mar  7 08:31 .gitignore
// -rw-r--r--   1 sungtaekim  staff    39B Mar  7 08:31 README.md
// drwxr-xr-x  28 sungtaekim  staff   896B Apr  4 09:16 cmd
// -rw-r--r--   1 sungtaekim  staff    23B Mar  7 08:31 golang_example.go

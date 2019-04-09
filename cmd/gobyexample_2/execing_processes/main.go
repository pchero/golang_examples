// https://gobyexample.com/execing-processes

// In the previous example we looked at
// spawning external processes. We
// do this when we need an external process accessible to
// a running Go process. Sometimes we just want to
// completely replace the current Go process with another
// (perhaps non-Go) one. To do this, we'll use Go's
// implementation of the classic
// exec
// function.

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// For our example, we'll exec 'ls'. Go requires an
	// absolute path to the binary we want to execute, so
	// we'll use 'exec.LookPath' to find it (probably
	// '/bin/ls').
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}

	// 'Exec' requires arguments in slice form (as
	// apposed to one big string). We'll give 'ls' a few
	// common arguments. Note that the first argument should
	// be the program name.
	args := []string{"ls", "-a", "-l", "-h"}

	// 'Exec' also needs a set of environment variables
	// to use. Here we just provides our current
	// environment.
	env := os.Environ()

	// Here's the actual 'syscall.Exec' call. If this call is
	// successful, the execution of our process will and
	// here and be replaced  by the '/bin/ls -a -l -h'
	// process. If there is an error we'll get a return
	// value.
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

// $ go run cmd/gobyexample_2/exiting_processes/main.go
// total 40
// drwxr-xr-x   8 sungtaekim  staff   256B Mar 14 09:17 .
// drwxr-xr-x   4 sungtaekim  staff   128B Mar 14 09:17 ..
// -rw-r--r--@  1 sungtaekim  staff   6.0K Mar 14 09:17 .DS_Store
// drwxr-xr-x  15 sungtaekim  staff   480B Apr  9 09:20 .git
// -rw-r--r--   1 sungtaekim  staff    17B Mar  7 08:31 .gitignore
// -rw-r--r--   1 sungtaekim  staff    39B Mar  7 08:31 README.md
// drwxr-xr-x  28 sungtaekim  staff   896B Apr  4 09:16 cmd
// -rw-r--r--   1 sungtaekim  staff    23B Mar  7 08:31 golang_example.go

// Note that Go does not offer a classic Unix fork function.
// Usually this isn't an issue though, since starting
// goroutines, spawning processes, and exec'ing processes
// covers most use cases for fork.

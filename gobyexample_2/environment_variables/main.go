// https://gobyexample.com/environment-variables

// Environment variables
// are a universal mechanism for conveying configuration
// information to Unix programs.
// Let's look at how to set, and list environment variables.

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// To set a key/value pair, use 'os.Setenv'. To get a
	// value for a key, use 'os.Getenv'. This will return
	// an empty string if the key isn't present in the
	// environment.
	os.Setenv("FOO", "1")
	fmt.Println("FOO: ", os.Getenv("FOO"))
	fmt.Println("BAR: ", os.Getenv("BAR"))

	// Use 'os.Environ' to list all key/value paris in the
	// environement. This returns a slice of strings in the
	// form 'KEY=value'. You can 'strings.Split' them to
	// get the key and value. Here we print all the keys.
	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0])
	}
}

// $ go run cmd/gobyexample_2/environment_variables/main.go
// FOO:  1
// BAR:

// SHELL
// LSCOLORS
// ITERM_PROFILE
// COLORTERM
// LESS
// XPC_FLAGS
// TERM_PROGRAM_VERSION
// SSH_AUTH_SOCK
// TERM_SESSION_ID
// GOOGLE_APPLICATION_CREDENTIALS_PERSONAL
// PWD
// LOGNAME
// GOOGLE_APPLICATION_CREDENTIALS
// COMMAND_MODE
// ITERM_SESSION_ID
// HOME
// LANG
// SECURITYSESSIONID
// TMPDIR
// TERM
// ZSH
// USER
// COLORFGBG
// SHLVL
// PAGER
// XPC_SERVICE_NAME
// LC_CTYPE
// Apple_PubSub_Socket_Render
// PATH
// OLDPWD
// GOPATH
// __CF_USER_TEXT_ENCODING
// TERM_PROGRAM
// APPLICATION_INSIGHTS_NO_DIAGNOSTIC_CHANNEL
// _
// FOO

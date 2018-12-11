package main

import (
	"flag"
	"io"
	"net/http"
	"os"
)

var (
	help     = flag.Bool("help", false, "Display uage information.")
	pemFile  = flag.String("pemfile", "./test.pem", "File path for pemfile.")
	username = flag.String("username", "test", "Allowed ssl username.")
	password = flag.String("password", "test", "Allowed ssl password.")
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/test", HandleIndex)

	// http.ListenAndServe(":8081", nil)
	http.ListenAndServeTLS(":8080", *pemFile, *pemFile, nil)

}

// HandleIndex default http callback
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	ret := basicAuth(r)
	if ret != true {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	io.WriteString(w, "hello, world\n")

}

// type handler func(w http.ResponseWriter, r *http.Request)

func basicAuth(r *http.Request) bool {

	// get authorization
	user, pass, ret := r.BasicAuth()
	if ret != true || *username != user || *password != pass {
		return false
	}

	return true
}

type testHandle int

func (n testHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ret := basicAuth(r)
	if ret != true {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	return
}

// func (c Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// }

// func basicAuthWrapper() http.Handler {

// 	return ServeHTTP(w, r) {

// 		return
// 	}

// 	// return func(w http.ResponseWriter, r *http.Request) {

// 	// 	// get authorization
// 	// 	user, pass, ret := r.BasicAuth()
// 	// 	if ret != true || *username != user || *password != pass {
// 	// 		http.Error(w, "authorization failed", http.StatusUnauthorized)
// 	// 		return
// 	// 	}

// 	// 	h(w, r)
// 	// 	return
// 	// }
// }

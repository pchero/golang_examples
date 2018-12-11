package main

import (
	"log"
	"net/http"
)

func main() {
	// public views
	http.HandleFunc("/", HandleIndex)

	// private views
	http.HandleFunc("/post", PostOnly(basicAuth(HandlePost)))
	http.HandleFunc("/json", GetOnly(basicAuth(HandleJSON)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// https://gowebexamples.com/json/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// User struct
type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func main() {
	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		fmt.Fprintf(w, "%s %s is %d yeard old!", user.Firstname, user.Lastname, user.Age)
	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		peter := User{
			Firstname: "John",
			Lastname:  "Doe",
			Age:       25,
		}

		json.NewEncoder(w).Encode(peter)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

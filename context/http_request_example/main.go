package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", StatusPage)
	mux.HandleFunc("/login", LoginPage)
	mux.HandleFunc("/logout", LogoutPage)

	log.Println("Start server on port :8085")
	log.Fatal(http.ListenAndServe(":8085", mux))
}

// StatusPage shows current status
func StatusPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This page will show the context username once the context is added."))
}

// LoginPage shows login page
func LoginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(365 * 24 * time.Hour) // set to expire in 1 year
	cookie := http.Cookie{
		Name:    "username",
		Value:   "pchero21@gmail.com",
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)
}

// LogoutPage shos logout page
func LogoutPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().AddDate(0, 0, -1) // set to expire in the past
	cookie := http.Cookie{
		Name:    "username",
		Value:   "pchero21@gmail.com",
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)
}

func AddContext(next http.Handler) http.Handler {
	return http.Handler(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "-", r.RequestURI)
		cookie, _ := r.Cookie("username")
		if cookie != nil {
			// Add data to context
			ctx := context.WithValue(r.Context(), "Username", cookie.Value)
			next.ServHTTP(w, r.WithContext(ctx))
		} else {
			next.ServHTTP(w, r)
		}
	})
}

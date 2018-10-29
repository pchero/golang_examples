// url provide a uniform way to locate resources
// here's how to parse URLs in go.

package main

import (
	"fmt"
	"net"
	"net/url"
)

func main() {
	// we'll parse this example URL, which includes a
	// scheme, authetication info, host, port, path,
	// query params, and query fragment.
	s := "postgres://user:pass@host.com:5432/path?k=v#f"

	// Parse the URL and ensure there are no errors.
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	// accessing the scheme is straitforward.
	fmt.Println(u.Scheme)

	// 'User' contains all authentication info; call
	// 'Username' and 'Password' on this for individual
	// values.
	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	p, _ := u.User.Password()
	fmt.Println(p)

	// the 'host' contains both the hostname and the port,
	// if present. use 'SplitHostPort' to extract them.
	fmt.Println(u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	fmt.Println(host)
	fmt.Println(port)

	// here we extract the 'path' and the fragment after
	// the '#'
	fmt.Println(u.Path)
	fmt.Println(u.Fragment)

	// to get query params in a string of 'k=v' format,
	// use 'RawQuery'. you can also parse query params
	// into a map. this parsed query param maps are from
	// strings to slices of strings, so index into '[0]'
	// if you only want the first value.
	fmt.Println(u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["k"][0])
}

package main

import (
	"fmt"
	"net/url"
	"path/filepath"
)

func main() {
	s := "postgres://user:pass@host.com:5432/path/test.wav?k=v#f"
	u, err := url.Parse(s)
	if err != nil {
		return
	}
	ext := filepath.Ext(u.Path)
	fmt.Printf("path: %s, ext: %s\n", u.Path, ext)

	s = "postgres://user:pass@host.com:5432/path/test?k=v#f"
	u, _ = url.Parse(s)
	ext = filepath.Ext(u.Path)
	fmt.Printf("path: %s, ext: %s\n", u.Path, ext)

	s = "postgres://user:pass@host.com:5432/path/test?k=v#f"
	u, _ = url.Parse(s)
	fmt.Printf("rawPath: %s, rawQuery: %s, string: %s\n", u.RawPath, u.RawQuery, u.String())
	// ext = filepath.Ext(u.RawPath)
	fmt.Printf("path: %s, ext: %s\n", u.Path, ext)
}

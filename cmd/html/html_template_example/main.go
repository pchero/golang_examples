package main

import (
	"fmt"
	"html"
	"html/template"
	"os"
)

func main() {
	const s = `"Fran & Freddie's Diner" <tasty@example.com>`
	const t = "&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;"
	v := []interface{}{`"Fran & Freddie's Diner"`, ' ', `<tasty@example.com>`}

	fmt.Println(template.HTMLEscapeString(s))
	template.HTMLEscape(os.Stdout, []byte(s))
	fmt.Fprintln(os.Stdout, "")

	fmt.Println(template.JSEscapeString(s))
	template.JSEscape(os.Stdout, []byte(s))
	fmt.Fprintln(os.Stdout, "")
	fmt.Println(template.JSEscaper(v...))

	fmt.Println(template.URLQueryEscaper(v...))

	if t != html.UnescapeString(t) {
		fmt.Println("It has HTML entities.")
		fmt.Println(html.UnescapeString(t))
	}

}

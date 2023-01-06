// https://gowebexamples.com/templates/

package main

import (
	"net/http"
	"text/template"
)

// Todo Default todo structure
type Todo struct {
	Title string
	Done  bool
}

// TodoPageData Default page structure
type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	tmpl := template.Must(template.ParseFiles("layout.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "My Todo list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}

		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}

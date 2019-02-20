// go's structs are typed collections of fields.
// they're useful for grouping data together to form
// records.

package main

import "fmt"

// this 'person' struct type has 'name' and 'age' fields.
type person struct {
	name string
	age  int
}

func main() {
	// this syntax creates a new struct.
	fmt.Println(person{"Bob", 20})

	// you can name the fields when initializing a struct.
	fmt.Println(person{name: "Alice", age: 30})

	// ommitted fields will be zero-valued
	fmt.Println(person{name: "Fred"})

	// an '&' prefix yields a pointer to the struct
	fmt.Println(&person{name: "Ann", age: 40})

	// access struct fields with a dot
	s := person{name: "Sean", age: 50}
	fmt.Println(s.name)

	// you can also use dots with struct pointers - the
	// pointers are automatically dereferenced.
	sp := &s
	fmt.Println(sp.age)

	// struct are mutable.
	sp.age = 51
	fmt.Println(sp.age)
}

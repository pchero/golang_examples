// switch statements express conditionals across many
// branches.

package main

import (
	"fmt"
	"time"
)

func main() {
	// here's a basic 'switch'.
	i := 2
	fmt.Println("Write ", i, " as ")
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	}

	// you can use commas to separate multiple expressions
	// in the same 'case' statement. we use the optional
	// 'default' case in this example as well.
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")
	case time.Monday:
		fallthrough
	default:
		fmt.Println("It's a weekday.")
	}

	// 'switch' without an expression is an alternate way
	// to express if/else logic. here we also show how the
	// 'case' expressions can be non-constants.
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}

	// a type 'switch' compares types instead of values. you
	// can use this to discover the type of an interface
	// value. in this example, the variable 't' will have the
	// type corresponding to its clause.
	whatAmI := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm an int")
		default:
			fmt.Printf("Don't know type %T\n", t)
		}
	}

	whatAmI(true)
	whatAmI(1)
	whatAmI("hey")
}

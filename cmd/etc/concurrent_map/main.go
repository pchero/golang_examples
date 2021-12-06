package main

import (
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type target struct {
	key string
	val int
}

func main() {
	m := cmap.New()
	// ch := make(chan int)
	const iterations = 100000000
	// var a [iterations]int

	t1 := &target{
		"k",
		0,
	}
	m.Set("test", t1)

	go func() {
		for i := 0; i < iterations; i++ {
			t1.val = i
		}
	}()

	go func() {
		tt, _ := m.Get("test")
		time.Sleep(time.Second)
		fmt.Printf("tt: %d\n", tt.(*target).val)
	}()

	go func() {
		time.Sleep(time.Microsecond * 500)
		m.Remove("test")
	}()

	// go func() {
	// 	for i := 0; i < iterations; i++ {
	// 		t1.val = i
	// 	}
	// }()
	// go func() {
	// 	for i := 0; i < iterations; i++ {
	// 		t1.val = i
	// 	}
	// }()
	// go func() {
	// 	for i := 0; i < iterations; i++ {
	// 		t1.val = i
	// 	}
	// }()

	// go func() {
	// 	for i := 0; i < iterations; i++ {
	// 		tmp, _ := m.Get("test")
	// 		if tmp.(*target).val != (iterations - 1) {
	// 			fmt.Printf("val: %d\n", tmp.(*target).val)
	// 		}
	// 	}
	// }()

	time.Sleep(time.Second * 10)

	// // Using go routines insert 1000 ints into our map.
	// go func() {
	// 	for i := 0; i < iterations/2; i++ {
	// 		// Add item to map.
	// 		m.Set(strconv.Itoa(i), i)

	// 		// Retrieve item from map.
	// 		val, _ := m.Get(strconv.Itoa(i))

	// 		// Write to channel inserted value.
	// 		ch <- val.(int)
	// 	} // Call go routine with current index.
	// }()

	// go func() {
	// 	for i := iterations / 2; i < iterations; i++ {
	// 		// Add item to map.
	// 		m.Set(strconv.Itoa(i), i)

	// 		// Retrieve item from map.
	// 		val, _ := m.Get(strconv.Itoa(i))

	// 		// Write to channel inserted value.
	// 		ch <- val.(int)
	// 	} // Call go routine with current index.
	// }()

	// // Wait for all go routines to finish.
	// counter := 0
	// for elem := range ch {
	// 	a[counter] = elem
	// 	counter++
	// 	if counter == iterations {
	// 		break
	// 	}
	// }

	// // Sorts array, will make is simpler to verify all inserted values we're returned.
	// sort.Ints(a[0:iterations])

	// // Make sure map contains 1000 elements.
	// if m.Count() != iterations {
	// 	logrus.Error("Expecting 1000 elements.")
	// }

	// // Make sure all inserted values we're fetched from map.
	// for i := 0; i < iterations; i++ {
	// 	if i != a[i] {
	// 		logrus.Errorf("missing value", i)
	// 	}
	// }
}

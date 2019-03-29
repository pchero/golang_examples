package main

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func main() {
	id := uuid.Must(uuid.NewV4())

	fmt.Printf("uuid info: %s\n", id)

	fmt.Printf("uuid info: %s\n", uuid.Must(uuid.NewV4()).String())
}

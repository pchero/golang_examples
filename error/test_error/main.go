package main

import (
	"fmt"
	"os"
)

func main() {
	// Attempt to open a file that doesn't exist
	file, err := os.Open("non_existent_file.txt")
	if err != nil {
		// Print the error message
		tmp := fmt.Errorf("emergencyAddressRepo.ListByParams.Builder: %w", err)
		fmt.Printf("Error: %w", tmp)
		return
	}
	defer file.Close()

	// If no error, proceed to read the file (this part won't be reached if there's an error)
	fmt.Println("File opened successfully")
}

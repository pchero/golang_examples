package main

import (
	"testoapi/gens/server" // Adjust this import path based on your project structure

	"github.com/gin-gonic/gin"
)

type Server struct{}

func strPtr(s string) *string {
	return &s
}

func main() {
	// Create the Gin app
	app := gin.Default()

	// Initialize the StrictServer implementation
	testServer := &Server{}

	v1 := app.RouterGroup.Group("v1.0")
	// // Create a new Strict handler using the generated NewStrictHandler function
	// handler := server.NewStrictHandler(strictServer, nil)

	// Register the routes using the generated RegisterHandlers function
	server.RegisterHandlers(v1, testServer)

	// Start the server on port 8080
	app.Run(":8080")
}

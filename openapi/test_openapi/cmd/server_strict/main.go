package main

// import (
// 	"context"
// 	"fmt"
// 	"testoapi/gencodes/models/customer_manager"
// 	"testoapi/gencodes/server" // Adjust this import path based on your project structure

// 	"github.com/gin-gonic/gin"
// )

// // StrictServerInterface Implementation
// type StrictServer struct{}

// // func (s *StrictServer) CreateAction(ctx context.Context, request server.CreateActionRequestObject) (server.CreateActionResponseObject, error) {
// // 	res := flow_manager.Action{
// // 		Id: "dd318f20-cda0-11ef-b7c1-5f9dbd04de67",
// // 	}
// // 	return server.CreateAction201JSONResponse(res), nil
// // }

// // // CreateAddress handler for the CreateAddress endpoint
// // func (s *StrictServer) CreateAddress(ctx context.Context, request server.CreateAddressRequestObject) (server.CreateAddressResponseObject, error) {
// // 	// Add your actual implementation here

// // 	tmpTarget := "test"
// // 	tmpType := common.AddressTypeAgent
// // 	res := common.Address{
// // 		Type:   &tmpType,
// // 		Target: &tmpTarget,
// // 	}

// // 	return server.CreateAddress201JSONResponse(res), nil
// // }

// func strPtr(s string) *string {
// 	return &s
// }

// // GetExample handler for the GetExample endpoint
// func (h *StrictServer) GetV10Accesskeys(ctx context.Context, request server.GetV10AccesskeysRequestObject) (server.GetV10AccesskeysResponseObject, error) {

// 	pageSize := 0
// 	if request.Params.PageSize != nil {
// 		pageSize = *request.Params.PageSize
// 	}

// 	pageToken := ""
// 	if request.Params.PageToken != nil {
// 		pageToken = *request.Params.PageToken
// 	}

// 	fmt.Printf("request detail. page_size: %d, page_token: %s, request: %v\n", pageSize, pageToken, request)

// 	result := []customer_manager.Accesskey{
// 		{
// 			Id: strPtr("619cfcd0-ce2e-11ef-8115-fba7cc388145"),
// 		},
// 	}

// 	// Add your actual implementation here
// 	return server.GetV10Accesskeys200JSONResponse{
// 		NextPageToken: strPtr("next_page_namge"),
// 		Result:        &result,
// 	}, nil
// }

// func (h *StrictServer) PostV10Accesskeys(ctx context.Context, request server.PostV10AccesskeysRequestObject) (server.PostV10AccesskeysResponseObject, error) {
// 	fmt.Printf("request: %v\n", request.Body)

// 	ID := "c75884ee-ce2a-11ef-9b36-93c4817072e0"
// 	tmp := customer_manager.Accesskey{
// 		Id: &ID,
// 	}
// 	// res := flow_manager.Action{
// 	// 	Id: "dd318f20-cda0-11ef-b7c1-5f9dbd04de67",
// 	// }
// 	return server.PostV10Accesskeys201JSONResponse(
// 		tmp,
// 	), nil
// }

// func (h *StrictServer) DeleteAccesskey(ctx context.Context, request server.DeleteAccesskeyRequestObject) (server.DeleteAccesskeyResponseObject, error) {

// 	return server.DeleteAccesskeyResponseObject(server.DeleteAccesskey204Response{}), nil
// }

// func (h *StrictServer) GetAccesskey(ctx context.Context, request server.GetAccesskeyRequestObject) (server.GetAccesskeyResponseObject, error) {
// 	return server.GetAccesskeyResponseObject(server.GetAccesskey200JSONResponse{}), nil

// }

// func (h *StrictServer) UpdateAccesskey(ctx context.Context, request server.UpdateAccesskeyRequestObject) (server.UpdateAccesskeyResponseObject, error) {
// 	return server.UpdateAccesskey200JSONResponse{}, nil
// }

// func main() {
// 	// Create the Gin router
// 	router := gin.Default()

// 	// Initialize the StrictServer implementation
// 	strictServer := &StrictServer{}

// 	// Create a new Strict handler using the generated NewStrictHandler function
// 	handler := server.NewStrictHandler(strictServer, nil)

// 	// Register the routes using the generated RegisterHandlers function
// 	server.RegisterHandlers(router, handler)

// 	// Start the server on port 8080
// 	router.Run(":8080")
// }

// // package main

// // import (
// // 	"testoapi/server" // Adjust this import based on your project structure

// // 	"github.com/gin-gonic/gin"
// // )

// // type Server struct{}

// // // CreateAddress is a simple handler implementation for the CreateAddress endpoint
// // func (s *Server) CreateAddress(c *gin.Context) {
// // 	// This is just a placeholder for actual logic
// // 	c.JSON(201, gin.H{"message": "Address created"})
// // }

// // // GetExample is a simple handler implementation for the GetExample endpoint
// // func (s *Server) GetExample(c *gin.Context) {
// // 	// This is just a placeholder for actual logic
// // 	c.JSON(200, gin.H{"message": "Example fetched"})
// // }

// // func main() {
// // 	// Create the Gin router
// // 	router := gin.Default()

// // 	// Create a new instance of your server handler
// // 	serverHandler := &Server{}

// // 	// Register the routes using the generated RegisterHandlers function
// // 	server.RegisterHandlers(router, serverHandler)

// // 	// Start the server
// // 	router.Run(":8080")
// // }

// // package main

// // import (
// // 	"context"
// // 	"fmt"
// // 	"log"
// // 	"testoapi/models"
// // 	"testoapi/server"

// // 	"github.com/gin-gonic/gin"
// // )

// // // Implementing the StrictServerInterface
// // type MyStrictServer struct{}

// // // CreateAddress implements the CreateAddress endpoint (strict interface)
// // func (s *MyStrictServer) CreateAddress(ctx context.Context, request server.CreateAddressRequestObject) (server.CreateAddressResponseObject, error) {
// // 	// Example: Create the address based on the request body
// // 	// Replace with actual address creation logic

// // 	fmt.Printf("Hello CreateAddress\n")

// // 	res := models.CommonAddress{
// // 		Type: models.CommonAddressTypeAgent,
// // 	}
// // 	return server.CreateAddress201JSONResponse(res), nil
// // }

// // // GetExample implements the GetExample endpoint (strict interface)
// // func (s *MyStrictServer) GetExample(ctx context.Context, request server.GetExampleRequestObject) (server.GetExampleResponseObject, error) {
// // 	// Example: Return a simple 200 OK response
// // 	return server.GetExample200Response{}, nil
// // }

// // func main() {
// // 	// Create a new Gin router
// // 	router := gin.Default()

// // 	// Create an instance of your strict server
// // 	myStrictServer := &MyStrictServer{}

// // 	// Create strict handler with middlewares (empty in this example)
// // 	handler := server.NewStrictHandler(myStrictServer, nil)

// // 	// Register the strict handler with the router
// // 	server.RegisterHandlers(router, handler)

// // 	// Start the server
// // 	log.Println("Starting server on :8080")
// // 	if err := router.Run(":8080"); err != nil {
// // 		log.Fatalf("Failed to start server: %v", err)
// // 	}
// // }

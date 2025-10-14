package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
)

func main() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable not set.")
	}

	// Replace with the SG Message ID you want to look up
	sgMessageID := "gAER9nwlRl67jn9geZDX_A"

	endpoint := fmt.Sprintf("/v3/messages/%s", sgMessageID)

	request := sendgrid.GetRequest(apiKey, endpoint, "https://api.sendgrid.com")
	request.Method = "GET"

	response, err := sendgrid.API(request)
	if err != nil {
		log.Println("API call failed:", err)
		return
	}

	fmt.Println("Status Code:", response.StatusCode)
	fmt.Println("Body:", response.Body)
	fmt.Println("Headers:", response.Headers)
}

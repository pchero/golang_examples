package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	// Retrieve the SendGrid API key from an environment variable.  This is best practice
	// for security.  Alternatively, you can hardcode it for testing purposes, but
	// NEVER commit that to a public repository.
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable not set.")
	}

	// Sender and recipient email addresses
	from := mail.NewEmail("Example User", "test@voipbin.net")      // Replace with your verified SendGrid sender email
	to := mail.NewEmail("Example Recipient", "pchero21@gmail.com") // Replace with recipient email

	// Email subject and content
	subject := "Sending with SendGrid is Fun"
	plainTextContent := "and easy to do anywhere, even with Go!"
	htmlContent := "<strong>and easy to do anywhere, even with Go!</strong>"

	// Create the email message
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Initialize the SendGrid client
	client := sendgrid.NewSendClient(apiKey)

	// Send the email
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return
	}

	// // Check the response from SendGrid
	// fmt.Println(response.StatusCode)
	// fmt.Println(response.Body)
	// fmt.Println(response.Headers)
	fmt.Printf("Response: %v\n", response)

	if response.StatusCode != 202 { // 202 means accepted for processing
		log.Printf("Email failed to send. Status code: %d", response.StatusCode)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

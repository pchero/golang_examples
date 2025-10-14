package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridEvent represents the structure of a SendGrid Event Webhook payload.
type SendGridEvent struct {
	Email       string `json:"email"`
	Timestamp   int64  `json:"timestamp"`
	Event       string `json:"event"`
	SgMessageId string `json:"sg_message_id"`

	MyMessageId string `json:"my_message_id,omitempty"` // Custom argument
	Status      string `json:"status,omitempty"`        // Bounce status code
	Reason      string `json:"reason,omitempty"`        // Bounce reason
}

// Function to get the email events
// Changed the parameter to sg_message_id to follow the official doc
func getEmailEvents(apiKey string, sgMessageID string) ([]SendGridEvent, error) {
	// Construct the API URL
	// The Message API uses "/messages/{msg_id}" to get the message
	baseURL := fmt.Sprintf("https://api.sendgrid.com/v3/messages/%s", sgMessageID)
	fmt.Printf("baseURL: %v\n", baseURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Add the API key to the request headers
	apiKeyEncoded := url.QueryEscape(apiKey) // Removing URL Encoding for now
	authHeader := "Bearer " + apiKeyEncoded
	fmt.Println("Authorization Header:", authHeader) // VERIFY HEADER
	req.Header.Add("Authorization", authHeader)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var event SendGridEvent
	err = json.Unmarshal(body, &event) // CHANGED THIS TO UNMARSHAL AN OBJECT INSTEAD OF ARRAY
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	return []SendGridEvent{event}, nil // Changed return type to SINGLE object

}

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

	// Create a V3 mail object
	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject

	// Create a plain text content and an html content
	message.AddContent(mail.NewContent("text/plain", plainTextContent))
	message.AddContent(mail.NewContent("text/html", htmlContent))

	// Create a personalization
	personalize := mail.NewPersonalization()
	personalize.AddTos(to)

	// Set the custom argument
	customMessageID := uuid.Must(uuid.NewV4()).String()
	fmt.Printf("customMessageID: %v\n", customMessageID)

	// customMessageID := "your-unique-message-id-123" // Replace with your actual ID
	personalize.SetCustomArg("my_message_id", customMessageID)

	// Add the personalization to the message
	message.AddPersonalizations(personalize)

	// Initialize the SendGrid client
	client := sendgrid.NewSendClient(apiKey)

	// Send the email
	resp, err := client.Send(message) // resp is *http.Response now

	if err != nil {
		log.Println(err)
		return
	}

	// // Check the response from SendGrid
	// fmt.Println(resp.StatusCode)
	// fmt.Println(resp.Body)
	// fmt.Println(resp.Headers)
	fmt.Printf("Response: %v\n", resp)
	// messageID := resp.Headers["X-Message-Id"][0]

	if resp.StatusCode != 202 { // 202 means accepted for processing
		log.Printf("Email failed to send. Status code: %d", resp.StatusCode)
	} else {
		fmt.Println("Email sent successfully!")
	}

	// Replace with the sg_message_id you want to query
	// Use the CustomMessageID set to the custom Args to find the SgMessageId
	// The retrieval is done using the SgMessageId

	sgMessageID := resp.Headers["X-Message-Id"][0]

	fmt.Printf("sgMessageID: %v\n", sgMessageID)

	// Retrieve email events from SendGrid API
	events, err := getEmailEvents(apiKey, sgMessageID)
	if err != nil {
		log.Fatalf("Error retrieving email events: %v", err)
		return
	}

	// Print out the information about the message that was sent
	if len(events) > 0 {
		for _, event := range events {
			fmt.Printf("Message ID: %s, Status: %s, Time: %s\n", event.SgMessageId, event.Event, time.Unix(event.Timestamp, 0))
			fmt.Printf("Message MyMessageId: %s, Status: %s, Time: %s\n", event.MyMessageId, event.Event, time.Unix(event.Timestamp, 0)) //check message
		}

	} else {
		fmt.Printf("no messages were sent with message ID %s", customMessageID)
	}

}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Telnyx API Endpoint
const telnyxBaseURL = "https://api.telnyx.com/v2"
const messagesEndpoint = "/messages"

// MessagePayload defines the structure of the JSON payload for sending an SMS.
type MessagePayload struct {
	From string `json:"from"`
	// To                 []string `json:"to"`
	To   string `json:"to"` // Changed to string for a single recipient
	Text string `json:"text"`
	// MessagingProfileID string   `json:"messaging_profile_id"`
	MediaURLs []string `json:"media_urls,omitempty"` // Optional
}

// MessageResponse represents the structure of the response from the Telnyx API.  Adjust this based on the actual response fields you need.
type MessageResponse struct {
	Data struct {
		ID         string    `json:"id"`
		RecordType string    `json:"record_type"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		From       string    `json:"from"`
		To         string    `json:"to"`
		Text       string    `json:"text"`
		Status     string    `json:"status"`
		Direction  string    `json:"direction"`
		Encoding   string    `json:"encoding"`
		Segments   int       `json:"segments"`
		Cost       string    `json:"cost"`
		Errors     []string  `json:"errors"`
		Media      []string  `json:"media"`
	} `json:"data"`
	Meta struct {
		Errors []struct {
			Code   string `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		} `json:"errors"`
	} `json:"meta"`
}

func main() {
	apiKey := os.Getenv("TELNYX_API_KEY")
	fromNumber := os.Getenv("TELNYX_SMS_PHONE_NUMBER")
	toNumber := os.Getenv("DESTINATION_PHONE_NUMBER")
	// toNumber := []string{"+821021656521", "+821021546521"}
	messagingProfileID := os.Getenv("TELNYX_SMS_PROFILE_ID")

	fmt.Printf("env info. apiKey: %v, fromNumber: %v, toNumber: %v, messaging_profile_id: %s\n", apiKey, fromNumber, toNumber, messagingProfileID)

	if apiKey == "" || fromNumber == "" || len(toNumber) == 0 {
		log.Fatal("Please set TELNYX_API_KEY, TELNYX_SMS_PHONE_NUMBER, and DESTINATION_PHONE_NUMBER environment variables.")
		return
	}

	message := MessagePayload{
		From: fromNumber,
		To:   toNumber,
		Text: "Hello from Telnyx raw API in Go! Let's go! Yeah!",
		// MessagingProfileID: messagingProfileID,
		MediaURLs: []string{}, // Optional: Add media URLs here
	}

	messageResponse, err := sendMessage(apiKey, message)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
		return
	}

	fmt.Printf("Message ID: %s\n", messageResponse.Data.ID)
	fmt.Printf("Message Status: %s\n", messageResponse.Data.Status)
	fmt.Printf("Message Created At: %s\n", messageResponse.Data.CreatedAt)
}

// sendMessage sends an SMS message using the Telnyx API.
func sendMessage(apiKey string, message MessagePayload) (*MessageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set a timeout
	defer cancel()

	// 1. Prepare the JSON payload
	payloadBytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	// 2. Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, telnyxBaseURL+messagesEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// 3. Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey) // Use Bearer token authentication

	// 4. Make the HTTP request
	client := &http.Client{} // You can customize the client (e.g., timeouts, transport)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// 5. Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	bodyString := string(bodyBytes)

	// 6. Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status code %d and body: %s", resp.StatusCode, bodyString)
	}

	// 7. Unmarshal the JSON response
	var messageResponse MessageResponse
	err = json.Unmarshal(bodyBytes, &messageResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w\nResponse Body: %s", err, bodyString)
	}

	if len(messageResponse.Meta.Errors) > 0 {
		return nil, fmt.Errorf("telnyx api returned errors: %+v", messageResponse.Meta.Errors)
	}

	return &messageResponse, nil
}

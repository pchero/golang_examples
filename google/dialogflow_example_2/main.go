package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/option"
)

// DecodeBase64 decodes a base64-encoded string
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// CreateDialogflowClient initializes a Dialogflow client using a base64-encoded credential
func CreateDialogflowClient(ctx context.Context, base64Cred string) (*dialogflow.SessionsClient, error) {
	decodedCred, err := DecodeBase64(base64Cred)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 credential: %v", err)
	}

	// Create a Dialogflow client with the credentials in-memory
	client, err := dialogflow.NewSessionsClient(ctx,
		option.WithCredentialsJSON(decodedCred),
		option.WithEndpoint("dialogflow.googleapis.com:443"),
		// option.WithEndpoint("europe-west1-dialogflow.googleapis.com:443"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Dialogflow client: %v", err)
	}

	return client, nil
}

// DetectIntent sends a text query to Dialogflow
func DetectIntent(ctx context.Context, client *dialogflow.SessionsClient, projectID, sessionID, query string) (string, error) {
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	req := &dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         query,
					LanguageCode: "en", // Change to your preferred language
				},
			},
		},
	}

	resp, err := client.DetectIntent(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to detect intent: %v", err)
	}
	fmt.Printf("Res: %v\n", resp) // Add newline for better readability.  Crucial for debugging

	// Extract and return the response text
	if resp.GetQueryResult() != nil {
		if resp.GetQueryResult().GetFulfillmentText() != "" { // Check if FulfillmentText is present
			return resp.GetQueryResult().GetFulfillmentText(), nil
		} else if len(resp.GetQueryResult().GetFulfillmentMessages()) > 0 { // Check for Fulfillment Messages

			for _, message := range resp.GetQueryResult().GetFulfillmentMessages() {
				if text := message.GetText(); text != nil && len(text.GetText()) > 0 {
					return text.GetText()[0], nil // Return the first text message, if any.
				}
			}
			return "", fmt.Errorf("no text fulfillment messages in Dialogflow response")
		} else {
			return "", fmt.Errorf("no fulfillment text or messages in Dialogflow response")
		}
	}

	return "", fmt.Errorf("no response from Dialogflow")
}

func main() {
	base64Cred := os.Getenv("CREDENTIAL_BASE64")
	projectID := os.Getenv("PROJECT_ID")

	ctx := context.Background()

	client, err := CreateDialogflowClient(ctx, base64Cred)
	if err != nil {
		log.Fatalf("Error creating Dialogflow client: %v", err)
	}
	defer client.Close()

	sessionID := "test-session1"
	query := "hi"

	response, err := DetectIntent(ctx, client, projectID, sessionID, query)
	if err != nil {
		log.Fatalf("Error detecting intent: %v", err)
	}

	fmt.Printf("Dialogflow Response: %s\n", response)
}

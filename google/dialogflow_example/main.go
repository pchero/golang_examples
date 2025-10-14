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
		option.WithEndpoint("europe-west1-dialogflow.googleapis.com:443"),
		// option.WithEndpoint("https://europe-west1-dialogflow.googleapis.com"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Dialogflow client: %v", err)
	}

	return client, nil
}

// DetectIntent sends a text query to Dialogflow
func DetectIntent(ctx context.Context, client *dialogflow.SessionsClient, projectID, sessionID, query string) (string, error) {
	// sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	// sessionPath := fmt.Sprintf("projects/%s/locations/%s/agents/%s/sessions/%s", projectID, "global", "Test-Chatbot-Date", sessionID)
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
	fmt.Printf("Res: %v", resp)

	// Extract and return the response text
	if resp.GetQueryResult() != nil {
		return resp.GetQueryResult().GetFulfillmentText(), nil
	}

	return "", fmt.Errorf("no response from Dialogflow")
}

func main() {

	base64Cred := os.Getenv("CREDENTIAL_BASE64")
	projectID := os.Getenv("PROJECT_ID")
	// Replace with your base64-encoded credentials and project ID
	// base64Cred := "YOUR_BASE64_CREDENTIAL"
	// projectID := "YOUR_PROJECT_ID"

	ctx := context.Background()

	client, err := CreateDialogflowClient(ctx, base64Cred)
	if err != nil {
		log.Fatalf("Error creating Dialogflow client: %v", err)
	}
	defer client.Close()

	sessionID := "test-session1"
	query := "what's flow?"

	response, err := DetectIntent(ctx, client, projectID, sessionID, query)
	if err != nil {
		log.Fatalf("Error detecting intent: %v", err)
	}

	fmt.Printf("Dialogflow Response: %s\n", response)
}

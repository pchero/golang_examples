package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

// AWS Credentials structure
type AWSCredentials struct {
	AccessKey string `json:"aws_access_key_id"`
	SecretKey string `json:"aws_secret_access_key"`
}

// AWS region map based on country codes
var awsRegions = map[string]string{
	"US": "us-east-1",
	"CA": "ca-central-1",
	"EU": "eu-central-1",
	"IN": "ap-south-1",
	"JP": "ap-northeast-1",
	"AU": "ap-southeast-2",
}

// GetRegionFromMetadata fetches the AWS region from EC2 instance metadata
func GetRegionFromMetadata() (string, error) {
	url := "http://169.254.169.254/latest/meta-data/placement/region"
	client := http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get region from AWS metadata: %v", err)
	}
	defer resp.Body.Close()

	region, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read region response: %v", err)
	}

	return string(region), nil
}

// GetPublicIP fetches the public IP of the machine
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("failed to get public IP: %v", err)
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read public IP response: %v", err)
	}

	return strings.TrimSpace(string(ip)), nil
}

// GetNearestAWSRegion gets the nearest AWS region based on public IP geolocation
func GetNearestAWSRegion() (string, error) {
	ip, err := GetPublicIP()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get geolocation: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		CountryCode string `json:"countryCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode geolocation response: %v", err)
	}

	// Map country to AWS region
	if region, found := awsRegions[result.CountryCode]; found {
		return region, nil
	}

	return "us-east-1", nil // Default fallback region
}

func main() {
	// Parse credential file flag
	credentialPath := flag.String("credential_file", "/aws_service_account.json", "AWS credential file path")
	flag.Parse()

	// Read AWS credentials from JSON file
	data, err := ioutil.ReadFile(*credentialPath)
	if err != nil {
		log.Fatalf("Failed to read credential file: %v", err)
	}

	var creds AWSCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		log.Fatalf("Failed to parse credential file: %v", err)
	}

	// // Get the nearest AWS region
	// region, err := GetRegionFromMetadata()
	// if err != nil {
	// 	log.Printf("Could not fetch AWS region from metadata, using geolocation: %v", err)
	// 	region, err = GetNearestAWSRegion()
	// 	if err != nil {
	// 		log.Fatalf("Failed to determine nearest AWS region: %v", err)
	// 	}
	// }

	region := "ap-northeast-2"

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(creds.AccessKey, creds.SecretKey, ""),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Initialize Polly client
	svc := polly.New(sess)
	fmt.Printf("Using AWS region: %s\n", region)

	// Text-to-Speech request
	inputText := "Hello, Welcome to the study world. This is a test message. Please enjoy the test world."
	input := &polly.SynthesizeSpeechInput{
		Text: aws.String(inputText),
		// OutputFormat: aws.String("mp3"),
		OutputFormat: aws.String(polly.OutputFormatPcm),
		VoiceId:      aws.String("Joanna"), // Change to another AWS Polly voice if needed
	}

	start := time.Now()

	// Generate speech
	resp, err := svc.SynthesizeSpeech(input)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	// Save audio to file
	filename := "output.wav"
	audioData, err := ioutil.ReadAll(resp.AudioStream)
	if err != nil {
		log.Fatalf("Failed to read audio stream: %v", err)
	}
	if err := ioutil.WriteFile(filename, audioData, 0644); err != nil {
		log.Fatalf("Failed to write audio file: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("SynthesizeSpeech took %s\n", elapsed)
	fmt.Printf("Audio content written to file: %v\n", filename)
}

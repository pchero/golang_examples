package main

import (
	"context"
	"fmt"
	"log"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

func main() {
	// Create a new Text-to-Speech client.
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// The text to be synthesized.
	text := "Thank you for calling. You've reached voicemail. Currently, the call cannot be answered. Please leave your message after the beep sound."

	// Perform the text-to-speech request on the text input with the selected voice parameters and audio file type.
	req := &texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		// Build the voice request; select the language code ("en-US") and the SSML voice gender ("FEMALE").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Wavenet-F",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:   texttospeechpb.AudioEncoding_LINEAR16,
			SampleRateHertz: 8000,
		},
	}

	// Perform the text-to-speech request.
	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	// The response's audioContent is binary.
	audioContent := resp.AudioContent
	// decodedAudio, err := base64.StdEncoding.DecodeString(string(audioContent))
	// if err != nil {
	// 	log.Fatalf("Failed to decode audio content: %v", err)
	// }

	// Write the binary audio content to a local file.
	err = os.WriteFile("voicemail.wav", audioContent, 0644)
	if err != nil {
		log.Fatalf("Failed to write audio content to file: %v", err)
	}

	fmt.Println("Audio content written to file: voicemail.wav")
}

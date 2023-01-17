package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var credentialPath = flag.String("credential_file", "./test.json", "GCP credential file path")

func main() {
	fmt.Printf("Hello world!\n")

	ctx := context.Background()
	// client, err := texttospeech.NewClient(ctx)
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(*credentialPath))
	if err != nil {
		logrus.Errorf("Could not create a new client. err: %v", err)
	}

	// perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: "Hello, Welcome to the study world. This is test message. Please enjoy the test world.",
			},
		},

		// build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral")
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			Name:         "en-US-Standard-C",
			// SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
			// SsmlGender: texttospeechpb.SsmlVoiceGender_MALE,
			SsmlGender: texttospeechpb.SsmlVoiceGender_FEMALE,
		},

		// select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		logrus.Fatal("Could not get correct response from the synthesize creator. err: %v", err)
	}

	// the resp's AudioContent is binary
	filename := "output.mp3"
	if err := ioutil.WriteFile(filename, resp.AudioContent, 0644); err != nil {
		logrus.Fatal("Could not crete result audio file. err: %v", err)
	}

	fmt.Printf("Audio content written to file: %v\n", filename)
}

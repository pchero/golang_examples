package main

import (
	"context"
	"fmt"
	"io/ioutil"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/sirupsen/logrus"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	fmt.Printf("Hello world!\n")

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		logrus.Errorf("Could not create a new client. err: %v", err)
	}

	// perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: "Hello, world!",
			},
		},

		// build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral")
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_NEUTRAL,
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

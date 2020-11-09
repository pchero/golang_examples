package main

import (
	"context"
	"flag"
	"io/ioutil"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var credentialPath = flag.String("credential_file", "./test.json", "GCP credential file path")

func main() {
	flag.Parse()

	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(*credentialPath))
	if err != nil {
		logrus.Errorf("Could not create a new client. err: %v", err)
		return
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
			AudioEncoding:   texttospeechpb.AudioEncoding_LINEAR16,
			SampleRateHertz: 8000,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		logrus.Errorf("Could not get a correct response. err: %v", err)
		return
	}

	filename := "output.wav"
	if err := ioutil.WriteFile(filename, resp.AudioContent, 0644); err != nil {
		logrus.Errorf("Could not create a result audio file. err: %v", err)
		return
	}

	return
}

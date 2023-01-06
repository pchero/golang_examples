package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/sirupsen/logrus"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

const usage = `Usage: main <audiofile>

Audio file must be a 16-bit signed listtle-endian encoded
with a sample rate of 16000.

The path to the audio file may be a GCS URI (gs://...).
`

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	var runFunc func(io.Writer, string) error
	path := os.Args[1]
	if strings.Contains(path, "://") {
		runFunc = recognizeGCS
	} else {
		runFunc = recognize
	}

	// perform the request
	if err := runFunc(os.Stdout, os.Args[1]); err != nil {
		logrus.Fatal(err)
	}

}

func recognizeGCS(w io.Writer, gcsURI string) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return err
	}

	// send the request with the URI(gs://...)
	// and sample rate information to be transcripted.
	// op, err := client.LongRunningRecognize(ctx, &speechpb.LongRunningRecognizeRequest{
	// 	Config: &speechpb.LongRunningRecognizeRequest{
	// 		// Encoding:        speechpb.RecognitionConfig_LINEAR16,
	// 		// SampleRateHertz: 16000,
	// 		// LanguageCode:    "en-US",
	// 		LanguageCode:    "en-US",
	// 		Encoding:        speechpb.RecognitionConfig_LINEAR16,
	// 		SampleRateHertz: 8000,
	// 	},
	// 	Audio: &speechpb.RecognitionAudio{
	// 		AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }

	// resp, err := op.Wait(ctx)
	// if err != nil {
	// 	return err
	// }

	// // print the results
	// for _, result := range resp.Results {
	// 	for _, alt := range result.Alternatives {
	// 		fmt.Fprintf(w, "\"%v\" (confidence=3f)\n", alt.Transcript, alt.Confidence)
	// 	}
	// }

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 8000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// print result
	fmt.Printf("Result: %v\n", resp)

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
	return nil
}

func recognize(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// send the contents of the audio file with the encoding and
	// sample rate information to be transcripted.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			// Encoding:        speechpb.RecognitionConfig_FLAC,
			// SampleRateHertz: 16000,
			LanguageCode:    "en-US",
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 8000,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v\n", resp)
	// print the result
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return nil
}

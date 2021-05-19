package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/sirupsen/logrus"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <AUDIOFILE>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<AUDIOFILE> must be a path to a local audio file. Audio file must be a 16-bit signed little-endian encoded with a sample rate of 16000.\n")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		logrus.Fatal("Please pass path to your local audio file as a command line argument")
	}
	audioFile := flag.Arg(0)
	logrus.SetLevel(logrus.DebugLevel)

	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	// stream, err := client.StreamingRecognize(ctx)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// // Send the initial configuration message.
	// if err := stream.Send(&speechpb.StreamingRecognizeRequest{
	// 	StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
	// 		StreamingConfig: &speechpb.StreamingRecognitionConfig{
	// 			Config: &speechpb.RecognitionConfig{
	// 				Encoding:        speechpb.RecognitionConfig_LINEAR16,
	// 				SampleRateHertz: 8000,
	// 				LanguageCode:    "en-US",
	// 			},
	// 		},
	// 	},
	// }); err != nil {
	// 	logrus.Fatal(err)
	// }

	f, err := os.Open(audioFile)
	if err != nil {
		logrus.Fatal(err)
	}
	defer f.Close()

	// go func() {
	buf := make([]byte, 48*1024) // 2byte(16bit) * 8000(hz) * 3 (seconds)
	// buf := make([]byte, 64)
	for {
		n, err := f.Read(buf)
		logrus.Debugf("Read byte: %d", n)
		if err == io.EOF {
			if _, err := os.Stat(audioFile); os.IsNotExist(err) {
				// finished
				break
			}

			logrus.Debugf("EOF. wait for 3 seconds.")
			time.Sleep(time.Second * 3)
			continue
		} else if err != nil {
			logrus.Errorf("Could not read the data from the file. err: %v", err)
			return
		}

		if n > 0 {
			logrus.Debugf("Sending a request.")
			resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 8000,
					LanguageCode:    "en-US",
				},
				Audio: &speechpb.RecognitionAudio{
					AudioSource: &speechpb.RecognitionAudio_Content{
						Content: buf[:n],
					},
				},
			})
			if err != nil {
				logrus.Errorf("Could not transcribe. err: %v", err)
				return
			}

			logrus.Debugf("Result: %v", resp)
			// Print the results.
			for _, result := range resp.Results {
				for _, alt := range result.Alternatives {
					logrus.Debugf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
				}
			}

			// if err := stream.Send(&speechpb.StreamingRecognizeRequest{
			// 	StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
			// 		AudioContent: buf[:n],
			// 	},
			// }); err != nil {
			// 	logrus.Printf("Could not send audio: %v", err)
			// }
			// }
			// if err == io.EOF {
			// 	// // Nothing else to pipe, close the stream.
			// 	// if err := stream.CloseSend(); err != nil {
			// 	// 	logrus.Fatalf("Could not close stream: %v", err)
			// 	// }
			// 	time.Sleep(time.Second * 3)
			// 	continue
			// } else if err != nil {
			// 	logrus.Printf("Could not read from %s: %v", audioFile, err)
			// 	return
			// }
		}
	}
	// }()

	// for {
	// 	resp, err := stream.Recv()
	// 	logrus.Debugf("Receiving stream. resp: %v", resp)

	// 	if err == io.EOF {
	// 		logrus.Debugf("Recevice EOF. err: %v", err)
	// 		time.Sleep(time.Second * 3)
	// 		continue
	// 	}
	// 	if err != nil {
	// 		logrus.Fatalf("Cannot stream results: %v", err)
	// 	}
	// 	if err := resp.Error; err != nil {
	// 		logrus.Fatalf("Could not recognize: %v", err)
	// 	}
	// 	for _, result := range resp.Results {
	// 		fmt.Printf("Result: %+v\n", result)
	// 	}
	// }
}

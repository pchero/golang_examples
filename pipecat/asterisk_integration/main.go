package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// 16bit mono PCM 샘플 wav 파일
	inputFile := "sample.wav"

	// GStreamer 파이프라인: wav 파일 → Opus → WebRTC
	pipeline := strings.Join([]string{
		"filesrc location=" + inputFile,
		"!",
		"wavparse",
		"!",
		"audioconvert",
		"!",
		"audioresample",
		"!",
		"opusenc",
		"!",
		"rtpopuspay",
		"!",
		"webrtcbin name=sendrecv",
	}, " ")

	args := strings.Fields(pipeline)

	cmd := exec.Command("gst-launch-1.0", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Starting GStreamer pipeline with file input...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run GStreamer pipeline: %v", err)
	}
}

// package main

// import (
// 	"bufio"
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"strings"

// 	"github.com/pion/webrtc/v3"
// )

// func main() {
// 	// 1. GStreamer 파이프라인 준비
// 	// pipeline := `
// 	//     appsrc name=src is-live=true format=time !
// 	//     audio/x-raw,format=S16LE,channels=1,rate=16000 !
// 	//     audioconvert ! audioresample !
// 	//     opusenc ! rtpopuspay !
// 	//     webrtcbin name=sendrecv
// 	// `
// 	pipeline := `-v appsrc name=src is-live=true format=time ! audio/x-raw,format=S16LE,channels=1,rate=16000 ! audioconvert ! audioresample ! opusenc ! rtpopuspay ! webrtcbin name=sendrecv`

// 	args := strings.Fields(pipeline)
// 	cmd := exec.Command("gst-launch-1.0", args...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Start(); err != nil {
// 		log.Fatalf("Failed to start GStreamer pipeline: %v", err)
// 	}
// 	log.Println("GStreamer pipeline started")

// 	// 2. Pion WebRTC PeerConnection 생성
// 	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// 3. 브라우저로부터 SDP Offer 받기
// 	fmt.Println("Paste SDP offer from browser (base64 encoded):")
// 	reader := bufio.NewReader(os.Stdin)
// 	offerB64, _ := reader.ReadString('\n')
// 	offerData, _ := base64.StdEncoding.DecodeString(strings.TrimSpace(offerB64))

// 	offer := webrtc.SessionDescription{
// 		Type: webrtc.SDPTypeOffer,
// 		SDP:  string(offerData),
// 	}

// 	if err := peerConnection.SetRemoteDescription(offer); err != nil {
// 		log.Fatal(err)
// 	}

// 	// 4. SDP Answer 생성 및 출력
// 	answer, err := peerConnection.CreateAnswer(nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := peerConnection.SetLocalDescription(answer); err != nil {
// 		log.Fatal(err)
// 	}

// 	answerB64 := base64.StdEncoding.EncodeToString([]byte(answer.SDP))
// 	fmt.Printf("SDP answer (paste to browser):\n%s\n", answerB64)

// 	// 5. stdin에서 slin16 PCM 데이터 읽어서 appsrc로 전달
// 	go func() {
// 		stdin := bufio.NewReader(os.Stdin)
// 		for {
// 			buf := make([]byte, 320) // 16kHz, 16bit, 10ms chunk
// 			n, err := stdin.Read(buf)
// 			if err != nil {
// 				log.Println("stdin read error:", err)
// 				break
// 			}
// 			if n > 0 {
// 				// TODO: Go에서 GStreamer appsrc로 push
// 				// 실제로는 gstreamer Go binding 필요
// 			}
// 		}
// 	}()

// 	select {} // 프로그램 종료 방지
// }

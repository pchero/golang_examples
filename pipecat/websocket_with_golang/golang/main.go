package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto" // Protobuf 직렬화/역직렬화를 위한 임포트

	pipecat_pb "test/models/pipecat"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// sendProtobufFrame은 Protobuf Frame을 직렬화하여 WebSocket으로 전송하는 헬퍼 함수입니다.
func sendProtobufFrame(ws *websocket.Conn, frame *pipecat_pb.Frame, clientAddr string) {
	log := logrus.WithFields(logrus.Fields{
		"func": "sendProtobufFrame",
	})

	marshaledFrame, err := proto.Marshal(frame)
	if err != nil {
		log.Printf("[WebSocket] Error marshaling Protobuf Frame for %s: %v", clientAddr, err)
		return
	}
	if err := ws.WriteMessage(websocket.BinaryMessage, marshaledFrame); err != nil {
		log.Printf("[WebSocket] Error writing Protobuf binary to client %s: %v", clientAddr, err)
		return
	}
}

func runTestPipecat(pipecatHandler PipecatcatHandler) error {
	log := logrus.WithFields(logrus.Fields{
		"func": "runTestPipecat",
	})

	filename, err := writeTestFile()
	if err != nil {
		return err
	}

	go func() {
		if errStart := pipecatHandler.Start("openai.gpt-4o-mini", "deepgram", "cartesia", filename); errStart != nil {
			log.Errorf("Error starting Pipecat handler: %v", errStart)
		}
	}()

	return nil
}

func stopTestPipecat(pipecatHandler PipecatcatHandler) error {
	pipecatHandler.Stop()

	return nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	log := logrus.WithFields(logrus.Fields{
		"func": "main",
	})
	log.Debugf("Hello world!")

	// write messages file

	pipecatHandler := NewPipecatHandler()

	if errRun := runTestPipecat(pipecatHandler); errRun != nil {
		log.Errorf("Failed to run test Pipecat: %v", errRun)
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	log.Debugf("Interrupt signal received.")
	log.Debugf("Terminating Python client...")
	stopTestPipecat(pipecatHandler)

	log.Debugf("Exiting main function")
}

func writeTestFile() (string, error) {
	log := logrus.WithFields(logrus.Fields{
		"func": "writeTestFile",
	})

	// Generate a new UUID
	// id := uuid.Must(uuid.NewV4())

	// Construct the filename: /tmp/<uuid>.json
	fileName := filepath.Join("/tmp", fmt.Sprintf("%s.json", uuid.Must(uuid.NewV4())))

	// The specific JSON data you provided
	fileData := []map[string]interface{}{
		{
			"role":    "system",
			"content": "You are a friendly AI assistant. When the conversation starts, always greet the user and ask how you can help. Respond naturally and keep your answers conversational. Start say 'hello' first.",
		},
	}

	// Marshal the Go slice of maps into a JSON byte slice
	jsonData, err := json.MarshalIndent(fileData, "", "  ") // Use MarshalIndent for pretty printing
	if err != nil {
		log.Errorf("Error marshaling JSON: %v", err)
		return "", err
	}

	// Write the JSON data to the file
	err = os.WriteFile(fileName, jsonData, 0644) // 0644 are the file permissions
	if err != nil {
		log.Errorf("Error writing file %s: %v", fileName, err)
		return "", err
	}

	log.Debugf("Successfully wrote JSON data to: %s", fileName)
	return fileName, nil
}

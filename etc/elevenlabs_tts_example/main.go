package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// TextToSpeechStream manages the WebSocket connection and audio processing.
type TextToSpeechStream struct {
	apiKey                   string
	voiceID                  string
	wsConn                   *websocket.Conn
	processor                AudioProcessor
	chanDone                 chan bool // Signals when handleWebSocketMessages goroutine has exited.
	chanFinalMessageReceived chan bool // Signals when the 'isFinal: true' message is received.
}

// Message struct to represent the message
type ElevenlabsMessage struct {
	Text                 string `json:"text"`
	TryTriggerGeneration bool   `json:"try_trigger_generation"`
	Finalize             bool   `json:"finalize"`
}

// AudioProcessor defines the interface for processing audio data.
type AudioProcessor interface {
	ProcessPCM(data []byte) error
	Close() error
}

// PCMFileWriter implements AudioProcessor to write PCM data to a file.
// It now includes sourceSampleRate (from API) and targetSampleRate (for file).
type PCMFileWriter struct {
	file             *os.File
	mu               sync.Mutex // Protects file writes.
	sourceSampleRate int        // Sample rate of incoming data from API (e.g., 16000)
	targetSampleRate int        // Desired sample rate for output file (e.g., 8000)
}

// NewPCMFileWriter creates a new PCMFileWriter instance.
// It now takes sourceSampleRate (expected from API) and targetSampleRate (desired for output file).
func NewPCMFileWriter(filename string, sourceSampleRate, targetSampleRate int) (*PCMFileWriter, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &PCMFileWriter{
		file:             f,
		sourceSampleRate: sourceSampleRate,
		targetSampleRate: targetSampleRate,
	}, nil
}

// ProcessPCM writes the given PCM data to the file, performing downsampling if necessary.
// This assumes 16-bit PCM data (2 bytes per sample).
func (w *PCMFileWriter) ProcessPCM(data []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.sourceSampleRate == w.targetSampleRate {
		// No resampling needed, write directly
		_, err := w.file.Write(data)
		return err
	} else if w.sourceSampleRate == 16000 && w.targetSampleRate == 8000 {
		// Ensure input data length is even for 16-bit samples
		if len(data)%2 != 0 {
			return fmt.Errorf("PCM data length must be even for 16-bit samples (received %d bytes)", len(data))
		}

		// Use a dynamically growing slice for downsampled data
		// Pre-allocate capacity for efficiency, roughly half the input size
		downsampledData := make([]byte, 0, len(data)/2)

		// Iterate through the input data, taking every other 16-bit sample.
		// Each iteration processes 4 bytes (two 16-bit samples) from the input
		// and appends 2 bytes (one 16-bit sample) to the output.
		for i := 0; i+3 < len(data); i += 4 { // Ensure at least 4 bytes (two 16-bit samples) are available for processing
			// Append the first 16-bit sample (data[i] and data[i+1])
			downsampledData = append(downsampledData, data[i], data[i+1])
		}

		_, err := w.file.Write(downsampledData)
		return err
	} else {
		return fmt.Errorf("unsupported sample rate conversion: from %d to %d", w.sourceSampleRate, w.targetSampleRate)
	}
}

// Close closes the underlying file.
func (w *PCMFileWriter) Close() error {
	return w.file.Close()
}

// NewTextToSpeechStream creates and initializes a new TextToSpeechStream.
func NewTextToSpeechStream(apiKey string, voiceID string, processor AudioProcessor) (*TextToSpeechStream, error) {
	// Construct the WebSocket URL for ElevenLabs.
	u := url.URL{
		Scheme: "wss",
		Host:   "api.elevenlabs.io",
		Path:   fmt.Sprintf("/v1/text-to-speech/%s/stream-input", voiceID),
	}

	// Add necessary query parameters for the stream.
	q := u.Query()
	q.Set("model_id", "eleven_multilingual_v2")
	// Set output_format to pcm_16000 as it's the minimum raw PCM output from ElevenLabs.
	q.Set("output_format", "pcm_16000")
	// Removed "sample_rate" query parameter as output_format dictates the output sample rate.
	u.RawQuery = q.Encode()

	// Set up WebSocket headers, including the API key.
	header := make(map[string][]string)
	header["xi-api-key"] = []string{apiKey}

	// Establish the WebSocket connection.
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return nil, fmt.Errorf("websocket connection failed: %v", err)
	}

	// Return a new TextToSpeechStream instance.
	return &TextToSpeechStream{
		apiKey:                   apiKey,
		voiceID:                  voiceID,
		wsConn:                   conn,
		processor:                processor,
		chanDone:                 make(chan bool),    // Unbuffered channel to signal goroutine completion.
		chanFinalMessageReceived: make(chan bool, 1), // Buffered channel to prevent deadlock if Stop is called before signal.
	}, nil
}

// AddText sends a text chunk to the ElevenLabs API for synthesis.
func (tts *TextToSpeechStream) AddText(text string) error {
	message := ElevenlabsMessage{
		Text:                 text,
		TryTriggerGeneration: true, // Suggests to the API to start generation if enough text is buffered.
	}

	return tts.wsConn.WriteJSON(message)
}

// FinalizeGeneration sends a message to indicate that no more text will be sent,
// prompting the API to send any remaining audio and the 'isFinal' message.
func (tts *TextToSpeechStream) FinalizeGeneration() error {
	message := ElevenlabsMessage{
		Text:                 "", // Empty text signifies the end of input.
		TryTriggerGeneration: true,
		Finalize:             true, // Explicitly tells the API to finalize generation.
	}
	return tts.wsConn.WriteJSON(message)
}

// Start begins listening for WebSocket messages in a new goroutine.
func (tts *TextToSpeechStream) Start() {
	go tts.handleWebSocketMessages()
}

// handleWebSocketMessages continuously reads messages from the WebSocket connection,
// processes audio data, and handles control messages.
func (tts *TextToSpeechStream) handleWebSocketMessages() {
	defer func() {
		// Close the WebSocket connection when the goroutine exits.
		// This defer will run whether the loop exits normally or due to an error.
		if err := tts.wsConn.Close(); err != nil {
			log.Printf("Error closing websocket connection in handleWebSocketMessages defer: %v", err)
		} else {
			log.Println("WebSocket connection closed by handleWebSocketMessages defer.")
		}
		// Signal that this goroutine is done.
		tts.chanDone <- true
		log.Println("handleWebSocketMessages goroutine signaled done.")
	}()

	for {
		messageType, message, err := tts.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket closed normally by server or client. Exiting handleWebSocketMessages.")
				// Signal finalMessageReceived because the stream has ended, even if no explicit isFinal.
				select {
				case tts.chanFinalMessageReceived <- true:
					log.Println("Signaled final message received due to normal WebSocket close.")
				default:
					// Non-blocking send: already signaled or not needed.
				}
			} else {
				log.Printf("Error reading websocket message: %v. Exiting handleWebSocketMessages.", err)
			}
			return // Exit the goroutine on any read error or connection close.
		}

		log.Printf("Received WebSocket message (type: %d, size: %d bytes)", messageType, len(message))

		var response map[string]interface{}
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("Error parsing response: %v. Message: %s", err, string(message))
			continue
		}

		// Process audio data if present.
		if audioDataB64, ok := response["audio"].(string); ok {
			decodedAudio, decodeErr := base64.StdEncoding.DecodeString(audioDataB64)
			if decodeErr != nil {
				log.Printf("Error decoding base64 audio data: %v", decodeErr)
				continue
			}
			log.Printf("Decoded audio chunk size before processing: %d bytes.", len(decodedAudio)) // Added log
			if err := tts.processor.ProcessPCM(decodedAudio); err != nil {
				log.Printf("Error processing PCM data: %v", err)
			}
			log.Printf("Processed audio chunk of size %d bytes.", len(decodedAudio))
		}

		// Check for the 'isFinal' flag, which indicates the end of audio generation.
		if isFinal, ok := response["isFinal"].(bool); ok { // Check if 'isFinal' key exists and is a boolean
			log.Printf("isFinal status received: %t", isFinal)
			if isFinal {
				log.Println("Received final message for current generation.")
				select {
				case tts.chanFinalMessageReceived <- true:
					log.Println("Signaled final message received.")
				default:
					log.Println("finalMessageReceived channel already signaled or not ready to receive.")
				}
			}
		}

		// Log other control messages like 'status'.
		if status, ok := response["status"].(string); ok {
			log.Printf("Status: %s", status)
		}

		// Log any errors reported by the server.
		if errMsg, ok := response["error"].(string); ok {
			log.Printf("Error from server: %s", errMsg)
		}
	}
}

// Stop sends a close message to the WebSocket and waits for the stream to finalize.
func (tts *TextToSpeechStream) Stop() {
	log.Println("Stop method called.")
	// Wait for the 'isFinal' message to be received from the server.
	// This ensures all audio data up to the final chunk has been processed.
	log.Println("Waiting for final message from ElevenLabs or timeout (30s)...")
	select {
	case <-tts.chanFinalMessageReceived:
		log.Println("Final message received. Proceeding to close stream.")
	case <-time.After(30 * time.Second): // Add a timeout to prevent indefinite waiting.
		log.Println("Timeout waiting for final message (30s). Closing stream anyway.")
	}

	// Send the 'close_stream' message to gracefully close the WebSocket connection.
	closeMsg := map[string]interface{}{
		"close_stream": true,
	}
	log.Println("Sending close_stream message...")
	// Check if the connection is still valid before attempting to write.
	// This prevents "use of closed network connection" error if the server already closed it.
	// Note: There's no direct IsOpen() for websocket.Conn, so we rely on the error.
	if err := tts.wsConn.WriteJSON(closeMsg); err != nil {
		log.Printf("Error sending close stream message: %v", err)
	} else {
		log.Println("close_stream message sent.")
	}

	// Wait for the handleWebSocketMessages goroutine to exit.
	// This ensures all cleanup in the goroutine is done, including closing the wsConn.
	log.Println("Waiting for handleWebSocketMessages goroutine to finish...")
	<-tts.chanDone
	log.Println("handleWebSocketMessages goroutine finished.")
	log.Println("TextToSpeechStream stopped successfully.")
}

func main() {
	// Retrieve the ElevenLabs API key from environment variables.
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		log.Fatal("ELEVENLABS_API_KEY environment variable is not set")
	}

	// Create a PCM file writer. We expect 16000 Hz from the API, but want 8000 Hz in the output file.
	pcmWriter, err := NewPCMFileWriter("output.raw", 16000, 8000)
	if err != nil {
		log.Fatal(err)
	}
	defer pcmWriter.Close() // Ensure the file is closed when main exits.

	// Define the voice ID to use for text-to-speech.
	voiceID := "21m00Tcm4TlvDq8ikWAM" // Default voice ID, replace with your preferred voice

	// Create a new TextToSpeechStream instance.
	tts, err := NewTextToSpeechStream(apiKey, voiceID, pcmWriter)
	if err != nil {
		log.Fatal(err)
	}

	// Start the WebSocket message handling goroutine.
	tts.Start()

	// Define the text chunks to be sent for synthesis.
	texts := []string{
		"Hello! This is a WebSocket test.",
		"We can add more text while processing.",
		"The audio quality is now set to 8 kilohertz PCM after downsampling.",
	}

	// Send each text chunk to the ElevenLabs API.
	for _, text := range texts {
		if err := tts.AddText(text); err != nil {
			log.Printf("Error adding text: %v", err)
		}
		// Small delay to simulate real-time text input.
		time.Sleep(100 * time.Millisecond)
	}

	// Add one more message after the initial batch.
	if err := tts.AddText("This is a delayed message!"); err != nil {
		log.Printf("Error adding delayed text: %v", err)
	}

	// Send the finalize message to tell ElevenLabs that no more text will be sent.
	// This prompts the API to send any remaining audio and the 'isFinal' message.
	if err := tts.FinalizeGeneration(); err != nil {
		log.Printf("Error finalizing generation: %v", err)
	}

	// Give the server a moment to send the final message after finalize.
	// Increased delay to provide more time for the 'isFinal' message.
	time.Sleep(2 * time.Second)

	// Call Stop to wait for the 'isFinal' message and then gracefully close the stream.
	tts.Stop()

	fmt.Println("Processing completed and stream closed!")
}

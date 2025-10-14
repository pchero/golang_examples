package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/CyCoreSystems/audiosocket"
	"github.com/google/uuid"
)

const (
	appName = "test"
	ariUser = "asterisk"
	ariPass = "asterisk"
	ariURL  = "http://localhost:8088/ari"
)

const (
	audioTypePCM16LE byte = 0x10 // signed linear, 16-bit, 8kHz, mono PCM (little-endian)
)

var (
	gBridgeHandle *ari.BridgeHandle
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen TCP port 8090 for audio socket connections
	listener, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		log.Fatalf("Failed to listen on TCP 8090: %v", err)
		return
	}
	defer listener.Close()

	go runListen(listener)

	// Connect to ARI
	client, err := native.Connect(&native.Options{
		Application: appName,
		Username:    ariUser,
		Password:    ariPass,
		URL:         ariURL,
	})
	if err != nil {
		log.Fatalf("Failed to connect to ARI: %v", err)
	}
	defer client.Close()

	if errEvent := runARIEventListen(ctx, client); errEvent != nil {
		log.Fatalf("Error running ARI event listener: %v", errEvent)
	}

	if errRun := runMainChannel(ctx, cancel, client); errRun != nil {
		log.Fatalf("Error running main channel: %v", errRun)
	}

	if errRun := runExternalMedia(ctx, cancel, client); errRun != nil {
		log.Printf("Error running external media: %v", errRun)
	}

	<-ctx.Done()
}

func runARIEventListen(ctx context.Context, client ari.Client) error {
	// Subscribe to ARI events
	sub := client.Bus().Subscribe(nil, "StasisStart", "StasisEnd", "ChannelDtmfReceived", "BridgeCreated")

	go func() {
		defer sub.Cancel()

		for {

			select {
			case e := <-sub.Events():
				log.Printf("-----------------------Received event: %T", e)

				switch ev := e.(type) {
				case *ari.StasisStart:
					ev, ok := e.(*ari.StasisStart)
					if !ok {
						log.Println("❌ Unexpected event type:", e)
						continue
					}
					log.Println("🔔 Channel started:", ev.Channel.ID)

					// Add the channel to the bridge
					err := gBridgeHandle.AddChannel(ev.Channel.ID)
					if err != nil {
						log.Printf("❌ Failed to add channel %s to bridge: %v", ev.Channel.ID, err)
					} else {
						log.Printf("✅ Channel added to bridge. channel_id: %s", ev.Channel.ID)
					}
					continue

				case *ari.StasisEnd:
					log.Println("🔔 Channel ended:", ev.Channel.ID)

				case *ari.ChannelDtmfReceived:
					log.Printf("🔔 DTMF received on channel %s: %s\n", ev.Channel.ID, ev.Digit)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func runMainChannel(ctx context.Context, cancel context.CancelFunc, client ari.Client) error {
	// Define bridge ID and key
	bridgeID := uuid.New().String()
	bridgeKey := ari.NewKey(ari.BridgeKey, bridgeID)

	// Create bridge
	var err error
	gBridgeHandle, err = client.Bridge().Create(bridgeKey, "mixing", "")
	if err != nil {
		log.Fatalf("Failed to create bridge: %v", err)
	}
	log.Println("✅ Bridge created:", gBridgeHandle.ID())

	log.Printf("test: %s", client.ApplicationName())

	// Originate call to PJSIP/1000
	tmp, err := client.Channel().Create(nil, ari.ChannelCreateRequest{
		Endpoint: "PJSIP/1000",
		App:      appName,
		AppArgs:  "hello world",
	})
	if err != nil {
		log.Fatalf("Failed to originate call: %v", err)
	}
	log.Printf("📞 Channel created to PJSIP/1000, channel ID: %s", tmp.ID())

	if errAdd := gBridgeHandle.AddChannel(tmp.ID()); errAdd != nil {
		log.Printf("❌ Failed to add channel %s to bridge: %v", tmp.ID(), errAdd)
	}
	log.Printf("✅ Channel %s added to bridge %s", tmp.ID(), gBridgeHandle)

	if errDial := tmp.Dial("", time.Second*30); errDial != nil {
		log.Printf("❌ Failed to dial channel %s: %v", tmp.ID(), errDial)
	} else {
		log.Printf("✅ Channel %s dialed successfully", tmp.ID())
	}

	return nil
}

func runListen(listener net.Listener) {
	fmt.Println("Listening for audio socket connections on port 8090...")

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Shutting down...")
		listener.Close()
		os.Exit(0)
	}()

	for {
		log.Printf("Waiting for AudioSocket connection...")
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Listener accept error: %v", err)
			continue
		}

		go handleAudioSocketRead(conn)
		go handleAudioSocketWrite(conn)
	}

}

func handleAudioSocketRead(conn net.Conn) {
	defer conn.Close()

	audiosocketID, err := audiosocket.GetID(conn)
	if err != nil {
		log.Printf("Failed to get AudioSocket ID: %v", err)
		return
	}
	log.Printf("Accepted AudioSocket connection with ID: %s", audiosocketID)

	// Keep the connection open until closed by remote (or add your own logic)
	buf := make([]byte, 1024)
	for {
		nRead, err := conn.Read(buf)
		if err != nil {
			break
		}
		log.Printf("Read %d bytes from AudioSocket connection.", nRead)
	}
}

func writeSampleData(conn net.Conn) error {
	decodedAudio, decodeErr := base64.StdEncoding.DecodeString(TestData)
	if decodeErr != nil {
		panic(decodeErr)
	}
	log.Printf("Decoded audio data length: %d bytes\n", len(decodedAudio))

	dataSamples, err := getDataSamples(16000, decodedAudio)
	if err != nil {
		panic(err)
	}
	log.Printf("Processed data length: %d bytes\n", len(dataSamples))

	// maxFragmentSize := 16000
	maxFragmentSize := 320
	offset := 0
	payloadLength := len(dataSamples)

	for offset < len(dataSamples) {
		fragmentLen := maxFragmentSize
		if offset+maxFragmentSize > payloadLength {
			fragmentLen = payloadLength - offset
		}

		fragment := dataSamples[offset : offset+fragmentLen]

		data, err := createAudioSocketPacket(audioTypePCM16LE, fragment)
		if err != nil {
			log.Printf("Failed to create audio socket packet: %v", err)
			return err
		}

		n, err := conn.Write(data)
		if err != nil {
			log.Printf("Failed to write to AudioSocket connection: %v", err)
			return err
		}
		log.Printf("Wrote %d bytes to AudioSocket connection. offset: %d, fragment length: %d", n, offset, fragmentLen)
		offset += fragmentLen

		time.Sleep(time.Millisecond * 20)
	}

	return nil
}

func handleAudioSocketWrite(conn net.Conn) {
	defer conn.Close()

	for {
		if errWrite := writeSampleData(conn); errWrite != nil {
			log.Printf("Failed to write sample data to AudioSocket connection: %v", errWrite)
			return
		}
		log.Printf("Sample data written to AudioSocket connection.")
	}
}

// getDataSamples processes 16-bit PCM data with the given inputRate sample rate.
// If inputRate equals defaultConvertSampleRate, it returns data as is.
// If inputRate is an integer multiple of defaultConvertSampleRate, it downsamples accordingly.
// Otherwise, it returns an error because only integer downsampling is supported.
func getDataSamples(inputRate int, data []byte) ([]byte, error) {
	if inputRate == 8000 {
		// No conversion needed
		return data, nil
	}

	if inputRate%8000 != 0 {
		return nil, fmt.Errorf("cannot convert %d Hz to %d Hz: only integer downsampling supported", inputRate, 8000)
	}

	factor := inputRate / 8000
	res := make([]byte, 0, len(data)/factor)

	// Downsample by selecting every 'factor'-th sample (2 bytes per sample)
	for i := 0; i+2*factor-1 < len(data); i += 2 * factor {
		res = append(res, data[i], data[i+1])
	}

	return res, nil
}

func createAudioSocketPacket(dataType byte, payload []byte) ([]byte, error) {
	// Create buffer
	buf := new(bytes.Buffer)

	// Write data type
	if err := buf.WriteByte(dataType); err != nil {
		return nil, fmt.Errorf("failed to write data type: %w", err)
	}

	// Write payload length (Big Endian)
	payloadLength := uint16(len(payload))
	if err := binary.Write(buf, binary.BigEndian, payloadLength); err != nil {
		return nil, fmt.Errorf("failed to write payload length: %w", err)
	}
	log.Printf("packet len: %d, data: %v\n", payloadLength, buf.Bytes())

	// Write payload
	_, err := buf.Write(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to write payload: %w", err)
	}

	return buf.Bytes(), nil
}

func runExternalMedia(ctx context.Context, cancel context.CancelFunc, client ari.Client) error {

	channelData, err := sendRequestExternalMediaChannel(
		context.Background(),
		"asterisk",
		"asterisk",
		appName,
		"127.0.0.1",
		8090,
	)
	if err != nil {
		log.Printf("failed to create external media channel: %v", err)
		return err
	}
	log.Printf("Created external media channel with ID: %s\n", channelData.ID)

	return nil
}

func sendRequestExternalMediaChannel(ctx context.Context, username, password, appName, externalHost string, externalMediaPort int) (*ari.ChannelData, error) {
	u, err := url.Parse("http://localhost:8088/ari/channels/externalMedia")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("app", appName)
	u.RawQuery = q.Encode()

	externalHostAddress := fmt.Sprintf("%s:%d", externalHost, externalMediaPort)
	audiosocketID := uuid.New().String()

	payload := map[string]any{
		"app":             appName,
		"channelId":       "844c7bb8-7683-11f0-8318-5bb53b1b2507",
		"external_host":   externalHostAddress,
		"encapsulation":   "audiosocket",
		"transport":       "tcp",
		"connection_type": "client",
		"format":          "slin",
		"direction":       "both",
		"data":            audiosocketID,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	log.Printf("Sending request to %s with body: %s", u.String(), string(bodyBytes))

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ARI returned non-success status: %s", resp.Status)
	}

	var channel ari.ChannelData
	if err := json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, fmt.Errorf("failed to decode ARI response: %w", err)
	}

	return &channel, nil
}

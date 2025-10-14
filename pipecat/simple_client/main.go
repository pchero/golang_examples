package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"unsafe"

	"github.com/gen2brain/malgo"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec"
	"github.com/pion/mediadevices/pkg/codec/g722" // ⬅️ FIX: nil pointer 패닉 우회를 위해 다시 추가 (go get 필요)
	_ "github.com/pion/mediadevices/pkg/driver/microphone"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v4"
)

// Pipecat signaling endpoint (adjust port/path if needed)
const pipecatOfferURL = "http://localhost:7860/api/offer"

// SDPOfferAnswer is used for JSON exchange with Pipecat.
type SDPOfferAnswer struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

// --- Global Audio Settings and Buffer (G.711 (PCMU) 설정) ---
var (
	audioBuffer                     = make([]byte, 0)
	bufferMutex                     sync.Mutex
	sampleRate                      = 8000
	playbackChannel                 = 1
	playbackFormat                  = malgo.FormatS16
	expectedRTPPayloadBytesPerFrame = 160
	expectedPCMOutputBytesPerFrame  = 320
)

// ulawToLinear는 µ-law 바이트를 16-bit linear PCM 샘플로 변환합니다. (표준 µ-law 디코더)
func ulawToLinear(u uint8) int16 {
	u = ^u
	sign := (u & 0x80)
	exponent := (u & 0x70) >> 4
	mantissa := u & 0x0F
	s := int(((int(mantissa) << 3) + 0x84) << exponent)
	s = s - 0x84
	if sign != 0 {
		s = -s
	}
	return int16(s)
}

// decodeULawPayloadToPCM은 µ-law 바이트 슬라이스를 PCM16 (Little Endian) 바이트 슬라이스로 변환합니다.
func decodeULawPayloadToPCM(payload []byte) []byte {
	out := make([]byte, len(payload)*2)
	for i, b := range payload {
		sample := ulawToLinear(b)
		binary.LittleEndian.PutUint16(out[i*2:], uint16(sample))
	}
	return out
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// WebRTC MediaEngine 설정
	m := &webrtc.MediaEngine{}

	// SDP 협상 범위를 G.711(PCMU) 및 PCMA에 맞춥니다. G722 등록을 제거했습니다.

	// 1. PCMU/8000 (Payload Type 0) 등록 - 최우선 코덱
	err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypePCMU,
			ClockRate: 8000,
			Channels:  uint16(playbackChannel),
		},
		PayloadType: 0,
	}, webrtc.RTPCodecTypeAudio)
	if err != nil {
		log.Fatal("PCMU 코덱 등록 오류:", err)
	}

	// 2. PCMA/8000 (Payload Type 8) 등록
	if err = m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:  webrtc.MimeTypePCMA,
			ClockRate: 8000,
			Channels:  uint16(playbackChannel),
		},
		PayloadType: 8,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		log.Fatal("PCMA 코덱 등록 오류:", err)
	}

	// 🛑 G722/8000 등록은 제거합니다. (G.711 클라이언트 목적에 충실)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	// PeerConnection 생성
	pc, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Fatal("PeerConnection 생성 오류:", err)
	}
	defer func() {
		if closeErr := pc.Close(); closeErr != nil {
			log.Println("PeerConnection 닫기 오류:", closeErr)
		}
	}()

	// PeerConnection 상태 변경 이벤트 처리
	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State 변경: %s\n", state.String())
		if state == webrtc.ICEConnectionStateFailed || state == webrtc.ICEConnectionStateClosed {
			fmt.Println("❌ WebRTC 연결이 실패하거나 닫혔습니다. 프로그램을 종료합니다.")
		}
	})

	// 수신 트랙 처리 핸들러 등록
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("🎶 원격 트랙 수신: ID %s, 종류 %s, 코덱: %s\n", track.ID(), track.Kind().String(), track.Codec().MimeType)

		if track.Kind() == webrtc.RTPCodecTypeAudio {
			go speakerOutput(track)
		}
	})

	// --- 마이크 설정 및 트랙 추가 ---

	// 🛑 FIX: G722를 사용하여 `pion/mediadevices`의 `nil pointer` 패닉을 우회합니다.
	// 이 인코더는 실제로 사용되지 않으며, `GetUserMedia`의 내부 구조를 초기화하는 용도입니다.
	g722Encoder, err := g722.NewEncoder(nil)
	if err != nil {
		// 사용자에게 go get 명령을 통해 모듈을 가져오도록 유도합니다.
		log.Fatal("G722 Encoder 생성 오류. 'go get github.com/pion/mediadevices/pkg/codec/g722'를 실행해야 할 수 있습니다:", err)
	}

	mediaStream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Audio: func(c *mediadevices.MediaTrackConstraints) {
			// 8kHz, 1채널로 마이크 입력 제약 설정
			c.SampleRate = prop.Int(sampleRate)
			c.ChannelCount = prop.Int(playbackChannel)
		},
		// CodecSelector 객체 초기화를 위해 더미 G722 코덱 명시
		Codec: mediadevices.CodecConfig{
			webrtc.RTPCodecTypeAudio: []codec.Codec{g722Encoder},
		},
	})
	if err != nil {
		log.Fatal("마이크 스트림 캡처 오류 (하드웨어/드라이버 문제일 수 있음):", err)
	}

	defer func() {
		for _, track := range mediaStream.GetTracks() {
			if closeErr := track.Close(); closeErr != nil {
				log.Printf("트랙 닫기 오류: %v", closeErr)
			}
		}
	}()

	for _, track := range mediaStream.GetTracks() {
		if track.Kind() == webrtc.RTPCodecTypeAudio {
			_, err = pc.AddTrack(track)
			if err != nil {
				log.Fatal("오디오 트랙 추가 오류:", err)
			}
			break
		}
	}

	// --- Offer 생성 및 전송 ---

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		log.Fatal("Offer 생성 오류:", err)
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		log.Fatal("LocalDescription 설정 오류:", err)
	}

	offerJSON := SDPOfferAnswer{
		SDP:  offer.SDP,
		Type: offer.Type.String(),
	}
	offerBody, err := json.Marshal(offerJSON)
	if err != nil {
		log.Fatal("Offer JSON 마샬링 오류:", err)
	}

	log.Println("➡️ Pipecat 서버로 Offer 전송...")
	log.Printf("💻 Offer SDP:\n%s", offer.SDP)

	resp, err := http.Post(pipecatOfferURL, "application/json", bytes.NewBuffer(offerBody))
	if err != nil {
		log.Fatalf("Pipecat 서버 POST 요청 오류 (URL: %s): %v", pipecatOfferURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Fatalf("Pipecat 서버 응답 오류: Status %d, Body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Answer JSON 수신 및 파싱
	var answerJSON SDPOfferAnswer
	if err := json.NewDecoder(resp.Body).Decode(&answerJSON); err != nil {
		log.Fatal("Answer JSON 디코딩 오류:", err)
	}

	log.Printf("💻 Answer SDP:\n%s", answerJSON.SDP)

	// Answer 수신 후 적용
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answerJSON.SDP,
	}
	// 이 SetRemoteDescription 호출에서 더 이상 패닉이 발생하지 않도록 수정되었습니다.
	err = pc.SetRemoteDescription(answer)
	if err != nil {
		log.Fatal("RemoteDescription 설정 오류:", err)
	}

	fmt.Println("✅ WebRTC 연결이 확립되었습니다. 마이크와 스피커가 활성화됩니다.")
	fmt.Println("---")
	fmt.Println("프로그램을 종료하려면 Ctrl+C를 누르십시오.")

	select {} // 무기한 대기
}

// speakerOutput decodes the RTP track and outputs it to the speaker.
func speakerOutput(track *webrtc.TrackRemote) {
	// 1. 디코더는 코덱 타입에 따라 동적으로 처리됩니다.
	decoderType := track.Codec().MimeType

	// 2. Malgo Context Initialization
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		log.Printf("Malgo context initialization error: %v", err)
		return
	}
	defer func() {
		if err := ctx.Uninit(); err != nil {
			log.Printf("Malgo context uninitialization error: %v", err)
		}
	}()

	// 3. Audio Device Setup and Initialization
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = playbackFormat
	deviceConfig.Playback.Channels = uint32(playbackChannel)
	deviceConfig.SampleRate = uint32(sampleRate)

	// Malgo Callback: Reads data from the buffer and sends it to the speaker
	onRecvFrames := func(pOutput, pInput []byte, framecount uint32) {
		framesToRead := int(framecount) * playbackChannel * int(unsafe.Sizeof(int16(0)))

		bufferMutex.Lock()
		defer bufferMutex.Unlock()

		if len(audioBuffer) >= framesToRead {
			copy(pOutput, audioBuffer[:framesToRead])
			audioBuffer = audioBuffer[framesToRead:]
		} else {
			for i := 0; i < len(pOutput); i++ {
				pOutput[i] = 0
			}
			log.Printf("⚠️ Audio buffer underflow! (Required: %d, Current: %d)", framesToRead, len(audioBuffer))
		}
	}

	var outputDevice *malgo.Device
	outputDevice, err = malgo.InitDevice(ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		log.Printf("Malgo device initialization error: %v", err)
		return
	}

	defer func() {
		outputDevice.Uninit()
	}()

	outputDevice.Start()

	// 4. RTP Packet Reception and Decoding Loop

	for {
		payload := make([]byte, 1500)
		payloadN, _, err := track.Read(payload)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("❌ RTP reception error: %v", err)
			break
		}

		var decodedBytes []byte

		switch decoderType {
		case webrtc.MimeTypePCMU:
			decodedBytes = decodeULawPayloadToPCM(payload[:payloadN])
		case webrtc.MimeTypePCMA:
			log.Printf("⚠️ PCMA 코덱이 협상되었으나, A-law 디코더는 구현되지 않았습니다. (PCMU 디코더만 지원)")
			continue
		default: // Opus, G722 등의 다른 모든 코덱 처리
			log.Printf("❌ 지원되지 않는 코덱 타입 (%s): PCMU 디코더만 지원합니다.", decoderType)
			continue
		}

		decodedBytesN := len(decodedBytes)

		if decodedBytesN > 0 {
			bufferMutex.Lock()
			audioBuffer = append(audioBuffer, decodedBytes...)
			bufferMutex.Unlock()
		} else {
			log.Printf("⚠️ G.711 디코딩 결과가 0 바이트입니다.")
		}
	}
}

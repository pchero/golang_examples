package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

const (
	pipecatURL  = "http://localhost:7860/api/offer"
	sampleRate  = 8000 // G.711은 8kHz 샘플링 레이트를 사용합니다.
	channels    = 1
	frameSizeMs = 20 // 20ms 오디오 프레임
	frameSize   = sampleRate * channels * frameSizeMs / 1000
)

func main() {
	// PortAudio 초기화
	if err := portaudio.Initialize(); err != nil {
		log.Fatalf("PortAudio 초기화 실패: %v", err)
	}
	defer portaudio.Terminate()

	// WebRTC MediaEngine 설정 및 G.711(PCMU) 코덱 등록
	m := &webrtc.MediaEngine{}
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypePCMU, ClockRate: sampleRate, Channels: channels},
		PayloadType:        0, // PCMU의 기본 페이로드 타입은 0
	}, webrtc.RTPCodecTypeAudio); err != nil {
		log.Fatalf("코덱 등록 실패: %v", err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	// PeerConnection 생성
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Fatalf("PeerConnection 생성 실패: %v", err)
	}
	defer peerConnection.Close()

	// <<< 추가: 데이터 채널 생성 및 이벤트 핸들러 설정
	// 클라이언트에서 'pipecat'이라는 이름으로 데이터 채널 생성
	dataChannel, err := peerConnection.CreateDataChannel("pipecat", nil)
	if err != nil {
		log.Fatalf("DataChannel 생성 실패: %v", err)
	}

	// 데이터 채널 상태 변경 핸들러
	dataChannel.OnOpen(func() {
		fmt.Println("WebRTC DataChannel 'pipecat'이 열렸습니다.")
		// 채널이 열리면 Pipecat에 준비 완료 메시지 등을 보낼 수 있습니다.
		// 이 메시지는 Pipecat의 설정에 따라 필요할 수도, 아닐 수도 있습니다.
		// 필요하다면, dataChannel.SendText("client ready") 등을 사용하세요.
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		// Pipecat으로부터 수신하는 애플리케이션 메시지 처리 (예: STT 결과, 턴 종료)
		fmt.Printf("DataChannel 메시지 수신: %s\n", string(msg.Data))
	})

	dataChannel.OnClose(func() {
		fmt.Println("WebRTC DataChannel 'pipecat'이 닫혔습니다.")
	})

	// 마이크 입력을 위한 오디오 트랙 생성
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypePCMU}, "audio", "pion")
	if err != nil {
		log.Fatalf("오디오 트랙 생성 실패: %v", err)
	}

	// PeerConnection에 오디오 트랙 추가
	rtpSender, err := peerConnection.AddTrack(audioTrack)
	if err != nil {
		log.Fatalf("트랙 추가 실패: %v", err)
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Pipecat으로부터 오디오 트랙을 수신했을 때의 핸들러
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("Pipecat으로부터 트랙 수신: %s\n", track.Codec().MimeType)

		out := make([]int16, frameSize)
		stream, err := portaudio.OpenDefaultStream(0, channels, float64(sampleRate), len(out), &out)
		if err != nil {
			log.Printf("PortAudio 출력 스트림 열기 실패: %v", err)
			return
		}
		defer stream.Close()

		if err := stream.Start(); err != nil {
			log.Printf("PortAudio 출력 스트림 시작 실패: %v", err)
			return
		}

		for {
			rtpPacket, _, err := track.ReadRTP()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Pipecat과의 오디오 스트림이 종료되었습니다.")
					return
				}
				log.Printf("RTP 패킷 읽기 오류: %v", err)
				return
			}

			payload := rtpPacket.Payload
			pcm := make([]int16, len(payload))
			for i, s := range payload {
				pcm[i] = ulawToLinear(s) // <<< CHANGED: G.711 u-law를 16-bit linear PCM으로 디코딩
			}

			if len(pcm) != len(out) {
				out = make([]int16, len(pcm))
			}

			copy(out, pcm)
			if err := stream.Write(); err != nil {
				log.Printf("오디오 출력 실패: %v", err)
			}
		}
	})

	// WebRTC 시그널링 (Offer/Answer 교환)
	// 1. Offer 생성
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Fatalf("Offer 생성 실패: %v", err)
	}

	// 2. LocalDescription 설정 및 ICE Candidate 수집 대기
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	if err := peerConnection.SetLocalDescription(offer); err != nil {
		log.Fatalf("LocalDescription 설정 실패: %v", err)
	}
	<-gatherComplete

	// <<< CHANGED: 생성된 Local Description (Offer) 로그 출력
	fmt.Println("------ [Client] SDP Offer ------")
	offerJSON, _ := json.MarshalIndent(peerConnection.LocalDescription(), "", "  ")
	fmt.Println(string(offerJSON))
	fmt.Println("--------------------------------")

	// 3. 생성된 Offer를 Pipecat 서버에 전송
	payload, err := json.Marshal(peerConnection.LocalDescription())
	if err != nil {
		log.Fatalf("Offer JSON 변환 실패: %v", err)
	}

	resp, err := http.Post(pipecatURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Pipecat 서버에 Offer 전송 실패: %v", err)
	}
	defer resp.Body.Close()

	// 4. Pipecat 서버로부터 Answer 수신
	var answer webrtc.SessionDescription
	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		log.Fatalf("Answer JSON 디코딩 실패: %v", err)
	}

	// <<< CHANGED: 수신한 Remote Description (Answer) 로그 출력
	fmt.Println("------ [Server] SDP Answer ------")
	answerJSON, _ := json.MarshalIndent(answer, "", "  ")
	fmt.Println(string(answerJSON))
	fmt.Println("---------------------------------")

	// 5. 수신한 Answer를 RemoteDescription으로 설정
	if err := peerConnection.SetRemoteDescription(answer); err != nil {
		log.Fatalf("RemoteDescription 설정 실패: %v", err)
	}

	fmt.Println("Pipecat과 WebRTC 연결이 성공적으로 수립되었습니다. 마이크 입력을 시작합니다.")

	// 마이크 입력을 처리하고 WebRTC 트랙으로 전송
	go func() {
		in := make([]int16, frameSize)
		stream, err := portaudio.OpenDefaultStream(channels, 0, float64(sampleRate), frameSize, &in)
		if err != nil {
			log.Fatalf("PortAudio 입력 스트림 열기 실패: %v", err)
		}
		defer stream.Close()

		if err := stream.Start(); err != nil {
			log.Fatalf("PortAudio 입력 스트림 시작 실패: %v", err)
		}

		for {
			if err := stream.Read(); err != nil {
				log.Printf("오디오 입력 읽기 실패: %v", err)
				continue
			}

			// 1. 묵음 데이터인지 확인
			// isSilence := checkSilence(in, SILENCE_THRESHOLD)

			data := linearToUlawPCM(in)

			// // 2. 로그 출력 (묵음일 때만 로그를 제한적으로 출력)
			// if isSilence {
			// 	log.Printf("Silent frame detected. ulaw_len: %d", len(data))
			// } else {
			// 	// 음성 데이터일 경우 상세 로그
			// 	log.Printf("Writing audio data to the webrtc track (VOICE). ulaw_len: %d", len(data))
			// }

			// data := linearToUlawPCM(in)
			// log.Printf("Writing audio data to the webrtc track. pcm_len: %d, ulaw_len: %d", len(in), len(data))

			if err := audioTrack.WriteSample(media.Sample{
				Data:     data, // <<< CHANGED: int16 PCM을 G.711(byte)로 인코딩
				Duration: time.Millisecond * frameSizeMs,
			}); err != nil {
				log.Printf("오디오 샘플 쓰기 실패: %v", err)
			}
		}
	}()

	// 프로그램이 종료되지 않도록 대기
	fmt.Println("Press Ctrl+C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("프로그램을 종료합니다.")
}

// linearToUlawPCM은 int16 PCM 슬라이스를 PCMU(G.711 u-law) 바이트 슬라이스로 변환합니다.
func linearToUlawPCM(samples []int16) []byte {
	muLawSamples := make([]byte, len(samples))
	for i, s := range samples {
		muLawSamples[i] = linearToUlaw(s)
	}
	return muLawSamples
}

// G.711 u-law 인코딩/디코딩 헬퍼 함수들
// (이전 답변과 동일하지만, 더 나은 음질을 위해 디코딩 함수 추가)
var (
	ulawExpLut = []int16{
		0, 132, 396, 924, 1980, 4092, 8316, 16764,
		-0, -132, -396, -924, -1980, -4092, -8316, -16764,
	}
)

// func linearToUlaw(pcmVal int16) byte {
// 	const (
// 		BIAS = 0x84
// 		MAX  = 32635
// 	)
// 	var mask uint16
// 	var seg uint16
// 	var uval byte

// 	pcm_val := pcmVal
// 	if pcm_val < 0 {
// 		pcm_val = BIAS - pcm_val
// 		mask = 0x7F
// 	} else {
// 		pcm_val += BIAS
// 		mask = 0xFF
// 	}
// 	if pcm_val > MAX {
// 		pcm_val = MAX
// 	}

// 	seg = 8
// 	var i int16 = 256
// 	for pcm_val < i {
// 		seg--
// 		i >>= 1
// 		if seg == 0 {
// 			break
// 		}
// 	}

// 	uval = byte((uint16(seg) << 4) | ((uint16(pcm_val) >> (uint16(seg) + 3)) & 0xF))
// 	if uval == 0 {
// 		uval = 2
// 	}

// 	return uval ^ byte(mask)
// }

func linearToUlaw(pcmVal int16) byte {
	const (
		BIAS = 0x84 // 132
		MAX  = 32635
	)
	var mask uint16
	var seg uint16

	pcm_val := pcmVal
	if pcm_val < 0 {
		pcm_val = BIAS - pcm_val
		mask = 0x7F
	} else {
		pcm_val += BIAS
		mask = 0xFF
	}
	if pcm_val > MAX {
		pcm_val = MAX
	}

	seg = 8
	for pcm_val < (1 << (seg + 4)) {
		seg--
	}

	uval := byte(((seg << 4) | ((uint16(pcm_val) >> (seg + 3)) & 0xF)) ^ mask)
	return uval
}

func ulawToLinear(ulaw byte) int16 {
	ulaw = ^ulaw
	mantissa := ((int16(ulaw) & 0x0F) << 3) + 0x84
	exponent := (int16(ulaw) & 0x70) >> 4
	pcm := (mantissa << exponent) - 0x84
	if (ulaw & 0x80) == 0 {
		return -pcm
	}
	return pcm
}

// checkSilence는 PCM 데이터의 음량을 측정하여 묵음인지 여부를 반환합니다.
func checkSilence(samples []int16, threshold int16) bool {
	var sumOfSquares int64
	for _, s := range samples {
		// RMS 계산을 위해 제곱의 합을 구합니다.
		// float 대신 int64를 사용하여 오버플로우를 방지하고 부동 소수점 연산을 피합니다.
		sumOfSquares += int64(s) * int64(s)
	}

	// RMS (Root Mean Square) 에너지 계산: sqrt(sum(x^2) / N)
	// 여기서는 sqrt까지 계산하지 않고, sumOfSquares / N 값만 threshold와 비교해도 됩니다.
	meanSquare := sumOfSquares / int64(len(samples))

	// threshold는 음량의 크기를 나타냅니다.
	// 이 값을 초과하면 음성, 아니면 묵음으로 간주합니다.
	return meanSquare < int64(threshold)*int64(threshold)
}

// 이 함수를 사용하기 위해, 메인 함수 상단에 임계값(Threshold)을 정의합니다.
const (
	// ... 기존 상수

	// 묵음으로 판단할 임계값 (Mean Square 값의 기준).
	// 이 값은 환경에 따라 조정해야 합니다 (예: 1000).
	// 실제 RMS 값으로 따지면 sqrt(1000*1000) = 1000 입니다.
	// 16bit PCM (-32768 ~ 32767) 범위에서 이 정도면 매우 낮은 음량입니다.
	SILENCE_THRESHOLD int16 = 1000
)

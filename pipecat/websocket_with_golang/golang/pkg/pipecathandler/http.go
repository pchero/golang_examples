package pipecathandler

import (
	"encoding/json"
	"net/http"
	"test/models/pipecat"
	pipecat_pb "test/models/pipecat"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *pipecatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithFields(logrus.Fields{
		"func": "ServeHTTP",
	})

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Could not upgrade to WebSocket: %v", err)
		return
	}
	defer func() {
		ws.Close()
		log.Debugf("WebSocket connection closed.")
	}()

	for {
		msgType, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Debugf("Client disconnected gracefully.")
			} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("Client disconnected unexpectedly: %v", err)
			} else {
				log.Errorf("Error reading message from client: %v", err)
			}
			break
		}

		switch msgType {
		case websocket.BinaryMessage:
			var frame pipecat_pb.Frame
			if errUnmarshal := proto.Unmarshal(message, &frame); errUnmarshal != nil {
				log.Errorf("Could not unmarshal Protobuf message: %v", errUnmarshal)
				continue
			}

			switch x := frame.Frame.(type) {
			case *pipecat_pb.Frame_Text:
				log.Debugf("Received TextFrame: ID=%d, Name=%s, Text='%s'", x.Text.Id, x.Text.Name, x.Text.Text)
				// responseFrame := &pipecat_pb.Frame{
				// 	Frame: &pipecat_pb.Frame_Text{
				// 		Text: &pipecat_pb.TextFrame{
				// 			Id:   x.Text.Id + 1, // ID를 증가시키는 예시
				// 			Name: "GoServerResponse",
				// 			Text: fmt.Sprintf("Go server received your text: '%s'", x.Text.Text),
				// 		},
				// 	},
				// }
				// if errSend := h.sendProtobufFrame(ws, responseFrame); errSend != nil {
				// 	log.Errorf("Could not send the frame.")
				// }

			case *pipecat_pb.Frame_Audio:
				log.Debugf("Received AudioFrame: ID=%d, Name=%s, SampleRate=%d, Channels=%d, AudioLen=%d")
				// log.Printf("[WebSocket] Received AudioRawFrame from %s: ID=%d, Name=%s, SampleRate=%d, Channels=%d, AudioLen=%d",
				// 	clientAddr, x.Audio.Id, x.Audio.Name, x.Audio.SampleRate, x.Audio.NumChannels, len(x.Audio.Audio))
				responseFrame := &pipecat_pb.Frame{
					Frame: &pipecat_pb.Frame_Audio{
						Audio: &pipecat_pb.AudioRawFrame{
							Id:          x.Audio.Id + 1,
							Name:        "GoServerAudioEcho",
							Audio:       x.Audio.Audio,
							SampleRate:  x.Audio.SampleRate,
							NumChannels: x.Audio.NumChannels,
							// Pts 필드는 optional이므로, 원래 값이 있으면 함께 보냅니다.
							Pts: x.Audio.Pts,
						},
					},
				}
				if errSend := h.sendProtobufFrame(ws, responseFrame); errSend != nil {
					log.Errorf("Could not send the frame.")
				}

			case *pipecat_pb.Frame_Transcription:
				log.Debugf("Received TranscriptionFrame: ID=%d, Name=%s, Text='%s', UserID=%s, Timestamp=%s")
				// responseFrame := &pipecat_pb.Frame{
				// 	Frame: &pipecat_pb.Frame_Transcription{
				// 		Transcription: &pipecat_pb.TranscriptionFrame{
				// 			Id:        x.Transcription.Id + 1,
				// 			Name:      "GoServerTranscriptionResponse",
				// 			Text:      fmt.Sprintf("Go server heard: '%s'", x.Transcription.Text),
				// 			UserId:    x.Transcription.UserId,
				// 			Timestamp: time.Now().Format(time.RFC3339),
				// 		},
				// 	},
				// }
				// if errSend := h.sendProtobufFrame(ws, responseFrame); errSend != nil {
				// 	log.Errorf("Could not send the frame.")
				// }

			case *pipecat_pb.Frame_Message:
				log.Debugf("Received MessageFrame: Data='%s'", x.Message.Data)
				if errMessage := h.receiveMessageFrameMessage([]byte(x.Message.Data)); errMessage != nil {
					log.Errorf("Could not process MessageFrame: %v", errMessage)
				}
				// MessageFrame에 대한 응답 (예시)
				// responseFrame := &pipecat_pb.Frame{
				// 	Frame: &pipecat_pb.Frame_Message{
				// 		Message: &pipecat_pb.MessageFrame{
				// 			Data: fmt.Sprintf("Go server received generic message: '%s'", x.Message.Data),
				// 		},
				// 	},
				// }
				// sendProtobufFrame(ws, responseFrame, clientAddr)

			default:
				log.Errorf("Could not recognize the Protobuf Frame type. type: %T", x)
			}

		case websocket.TextMessage:
			// because of we switched to Protobuf communication, text messages are only handled in exceptional cases.
			log.Errorf("Could not recognize the message type. type: %d (Expecting Protobuf Binary)", msgType)

		case websocket.CloseMessage:
			log.Debugf("Received Close message from client.")
			return
		case websocket.PingMessage:
			log.Debugf("Received Ping message from client. Sending Pong.")
			if errWrite := ws.WriteMessage(websocket.PongMessage, []byte{}); errWrite != nil {
				log.Errorf("Could not send Pong message: %v", errWrite)
				return
			}
		case websocket.PongMessage:
			log.Debugf("Received Pong message from client.")
		default:
			log.Debugf("Received unknown message type %d", msgType)
		}
	}
}

func (h *pipecatHandler) sendProtobufFrame(ws *websocket.Conn, frame *pipecat_pb.Frame) error {
	marshaledFrame, err := proto.Marshal(frame)
	if err != nil {
		return errors.Wrapf(err, "could not marshaling the protobuf frame")
	}
	if err := ws.WriteMessage(websocket.BinaryMessage, marshaledFrame); err != nil {
		return errors.Wrapf(err, "could not write message")
	}

	return nil
}

func (h *pipecatHandler) receiveMessageFrameMessage(message []byte) error {
	log := logrus.WithFields(logrus.Fields{
		"func": "receiveMessageFrameMessage",
	})

	frame := pipecat.CommonFrameMessage{}
	if errUnmarshal := json.Unmarshal(message, &frame); errUnmarshal != nil {
		log.Errorf("Error unmarshaling JSON message: %v", errUnmarshal)
		return errUnmarshal
	}
	log.Debugf("Received message frame: Label=%s, Type=%s", frame.Label, frame.Type)

	if frame.Label != pipecat.RTVIMessageLabel {
		// other message types can be handled here
		log.Errorf("Unrecognized message label: %s", frame.Label)
		return nil
	}

	switch frame.Type {
	default:
		log.Errorf("Unrecognized RTVI message type: %s", frame.Type)
	}

	return nil
}

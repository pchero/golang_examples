package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// engineOpenaiHandler는 OpenAI API 클라이언트를 캡슐화합니다.
type engineOpenaiHandler struct {
	client *openai.Client
}

// NewEngineOpenaiHandler는 engineOpenaiHandler의 새 인스턴스를 생성합니다.
func NewEngineOpenaiHandler(openaiAPIKey string) *engineOpenaiHandler {
	client := openai.NewClient(openaiAPIKey)
	return &engineOpenaiHandler{
		client: client,
	}
}

// sendStreaming은 OpenAI ChatCompletion 요청을 스트리밍 모드로 보내고,
// 모델로부터 받은 각 문장을 문자열 채널을 통해 반환합니다.
func (h *engineOpenaiHandler) sendStreaming(ctx context.Context, req *openai.ChatCompletionRequest) (<-chan string, error) {
	// 스트리밍 요청을 위해 Stream: true 설정
	req.Stream = true

	stream, err := h.client.CreateChatCompletionStream(ctx, *req)
	if err != nil {
		return nil, fmt.Errorf("CreateChatCompletionStream error: %w", err)
	}

	// 스트리밍된 토큰들을 전달할 채널
	outputChan := make(chan string)

	go func() {
		defer stream.Close()    // 스트림 사용이 끝나면 닫습니다.
		defer close(outputChan) // 채널 사용이 끝나면 닫습니다.

		var currentSentence strings.Builder

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				// 스트림 종료, 남은 문장이 있다면 전송
				if currentSentence.Len() > 0 {
					outputChan <- strings.TrimSpace(currentSentence.String())
				}
				break
			}
			if err != nil {
				// 스트림 에러 처리
				log.Printf("Stream error: %v", err)
				break
			}

			log.Printf("Received chunk: %+v", response.Choices)

			// 첫 번째 Choice만 처리 (대부분의 경우 하나만 옴)
			for _, choice := range response.Choices {
				if choice.Delta.Content != "" {
					currentSentence.WriteString(choice.Delta.Content)

					// 문장 끝을 나타내는 문자를 찾습니다 (마침표, 물음표, 느낌표, 줄바꿈)
					// 여기에 좀 더 정교한 문장 분리 로직을 넣을 수 있습니다.
					if strings.ContainsAny(choice.Delta.Content, ".?!\n") && currentSentence.Len() > 0 {
						sentence := currentSentence.String()
						trimmedSentence := strings.TrimSpace(sentence)

						if trimmedSentence != "" {
							outputChan <- trimmedSentence // 사용자에게 전달
						}
						currentSentence.Reset() // 버퍼 초기화
					}
				}
			}
		}
	}()

	return outputChan, nil
}

func main() {
	// 환경 변수에서 OpenAI API 키를 가져옵니다.
	// `export OPENAI_API_KEY="your_api_key_here"` 로 설정하거나 직접 입력하세요.
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		log.Fatal("OPENAI_API_KEY 환경 변수가 설정되지 않았습니다.")
	}

	handler := NewEngineOpenaiHandler(openaiAPIKey)
	ctx := context.Background()

	// OpenAI 요청 설정
	chatReq := &openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo, // 또는 openai.GPT4 등 사용 가능한 모델
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "한국의 수도에 대해 자세히 설명해줘. 역사, 지리적 특징, 문화적 중요성을 포함해서 길게 설명해줘.",
			},
		},
		MaxTokens: 500, // 최대 응답 토큰 수 설정 (선택 사항)
	}

	fmt.Println("OpenAI 스트리밍 응답 시작...")

	// sendStreaming 함수 호출
	outputChan, err := handler.sendStreaming(ctx, chatReq)
	if err != nil {
		log.Fatalf("Failed to start streaming: %v", err)
	}

	// 채널에서 받은 문장을 콘솔에 출력합니다.
	for sentence := range outputChan {
		fmt.Printf("> %s\n", sentence) // 각 문장 앞에 "> "를 붙여 구분
	}

	fmt.Println("모든 스트리밍 응답 수신 완료.")
}

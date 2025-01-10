package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// NATS 서버에 연결
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// JetStream 컨텍스트 생성
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// 스트림 생성
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "DELAYED",
		Subjects: []string{"delayed.*"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 메시지 발행 (10초 지연)
	msg := nats.NewMsg("delayed.message")
	msg.Data = []byte("이 메시지는 10초 후에 전달됩니다.")
	msg.Header.Set("Nats-Delay", "10s")

	_, err = js.PublishMsg(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("메시지가 발행되었습니다.")

	// 소비자 생성 및 메시지 처리
	sub, err := js.PullSubscribe("delayed.message", "DELAYED_CONSUMER")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(12 * time.Second)

	msgs, err := sub.Fetch(1)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range msgs {
		log.Printf("받은 메시지: %s\n", string(msg.Data))
		msg.Ack()
	}
}

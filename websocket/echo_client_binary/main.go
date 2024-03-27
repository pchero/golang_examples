package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	endpoint := fmt.Sprintf("wss://%s", *addr)
	log.Debugf("Connecting to the server. endpoint: %s", endpoint)

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Errorf("Could not parse the endpoint. err: %v", err)
		return
	}

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Errorf("Could not connect. err: %v", err)
	}
	defer c.Close()
	log.Debugf("Connected to the websocket")

	// p := &rtp.Packet{}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Debugf("read: %v", err)
			return
		}
		log.Debugf("Received message. len: %d", len(message))

		// // Unmarshal the packet and update the PayloadType
		// if errUnmarshal := p.Unmarshal(message); errUnmarshal != nil {
		// 	log.Errorf("Could not unmarshal the received data. len: %d", len(message))
		// 	continue
		// }
		// log.Debugf("received message. rtp: %v", p)

		// // check the payload type
		// if p.PayloadType > 63 && p.PayloadType < 96 {
		// 	// rtcp packet.
		// 	// we don't deal with the rtcp
		// 	continue
		// }

		if errWrite := c.WriteMessage(websocket.BinaryMessage, message); errWrite != nil {
			log.Errorf("Could not write message: %v", errWrite)
		}
	}
}

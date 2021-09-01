package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {
	listenIP := flag.String("ip", "", "listen ip")
	listenPort := flag.Int("port", 0, "listen port")
	flag.Parse()

	listenAddr := fmt.Sprintf("%s:%d", *listenIP, *listenPort)
	fmt.Printf("Listen address: %s\n", listenAddr)

	// conn, err := net.ListenUDP("udp", &net.UDPAddr{
	// 	IP:   net.ParseIP(*listenIP),
	// 	Port: *listenPort,
	// })
	conn, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 4096)
		rlen, remote, err := conn.ReadFrom(message[:])
		if err != nil {
			panic(err)
		}

		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received: %s from %s\n", data, remote)
	}
}

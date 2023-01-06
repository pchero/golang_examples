package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/sirupsen/logrus"
	"github.com/youpy/go-wav"
)

const maxRecvRTCPBuf = 1500   // max rtcp recv buf size
const defaultMTU = 1200       // udp mtu
const ulawSamplingRate = 8000 // ulaw sampling rate
const payloadTypePCMU = 0     // ulaw rtp packet type. see detail https://en.wikipedia.org/wiki/RTP_payload_formats

func main() {
	infile_path := flag.String("infile", "", "wav file to read")
	destIP := flag.String("ip", "", "destination ip")
	destPort := flag.Int("port", 0, "destination port")
	flag.Parse()

	file, _ := os.Open(*infile_path)
	reader := wav.NewReader(file)

	defer file.Close()

	if err := sendMedia(*destIP, *destPort, reader); err != nil {
		logrus.Errorf("Could not send the media. err: %v", err)
	}
}

// SendMedia sends the media rtp to the destination
func sendMedia(host string, port int, wavAudio *wav.Reader) error {
	log := logrus.WithField("func", "sendMedia")

	// init
	conn, err := sendMediaInit(host, port)
	if err != nil {
		log.Errorf("failed sendMediaInit", "error", err)
		return err
	}

	// send
	if errMediaHandle := sendMediaHandle(conn, wavAudio); errMediaHandle != nil {
		log.Errorf("Could not send the media correctly", "error", errMediaHandle)
		return errMediaHandle
	}

	return nil
}

// sendMediaInit inits the settings for the media send.
func sendMediaInit(host string, port int) (*net.UDPConn, error) {
	log := logrus.WithField("func", "sendMediaInit")

	address := fmt.Sprintf("%s:%d", host, port)
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Errorf("Could not resolve the remote address. address: %v, err: %v", address, err)
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Errorf("Could not connect to the remote. err: %v", err)
		return nil, err
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, maxRecvRTCPBuf)
		for {
			if _, _, rtcpErr := conn.ReadFromUDP(rtcpBuf); rtcpErr != nil {
				return
			}
			log.Debugf("Recevied RTCP. buf: %v", rtcpBuf)
		}
	}()

	return conn, nil
}

// sendMediaHandle handles the RTP media sending and RTCP receiving
func sendMediaHandle(conn *net.UDPConn, wavAudio *wav.Reader) error { //nolint:interfacer // this is ok
	log := logrus.WithField("func", "sendMediaHandle")

	// set rtp options(g711 ulaw)
	seq := rtp.NewRandomSequencer()
	ssrc := rand.Uint32()
	payloader := &codecs.G711Payloader{}
	packetrizer := rtp.NewPacketizer(defaultMTU, payloadTypePCMU, ssrc, payloader, seq, ulawSamplingRate)

	data := make([]byte, 4096)

	// 1/8000 = 125 mciroseconds
	// data byte * 125
	tickDuration := time.Microsecond * time.Duration(len(data)*125)
	ticker := time.NewTicker(tickDuration)
	for ; true; <-ticker.C {
		l, err := wavAudio.Read(data)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Errorf("Could not read the sample. err: %v", err)
			return err
		}

		// packetize to the RTP
		packets := packetrizer.Packetize(data, uint32(l))
		for _, packet := range packets {
			b, err := packet.Marshal()
			if err != nil {
				continue
			}

			// send
			_, err = conn.Write(b)
			if err != nil {
				log.Errorf("Could not send the rtp correctly. err: %v", err)
				return err
			}
		}
	}

	return nil
}

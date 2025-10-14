package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	audioTypePCM16LE byte = 0x10 // signed linear, 16-bit, 8kHz, mono PCM (little-endian)
)

func main() {
	// Example PCM data (10 samples, all zero)
	pcmData := []int16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// Convert PCM data to little-endian byte array
	var pcmBytes bytes.Buffer
	if err := binary.Write(&pcmBytes, binary.LittleEndian, pcmData); err != nil {
		log.Fatalf("Failed to write PCM data: %v", err)
	}

	// Create AudioSocket packet
	packet, err := createAudioSocketPacket(audioTypePCM16LE, pcmBytes.Bytes())
	if err != nil {
		log.Fatalf("Failed to create AudioSocket packet: %v", err)
	}

	fmt.Printf("AudioSocket Packet: %X\n", packet)
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

	// // Write payload
	// _, err := buf.Write(payload)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to write payload: %w", err)
	// }

	return buf.Bytes(), nil
}

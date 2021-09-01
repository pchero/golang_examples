package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/youpy/go-wav"
)

func main() {
	infile_path := flag.String("infile", "", "wav file to read")
	flag.Parse()

	file, _ := os.Open(*infile_path)
	reader := wav.NewReader(file)

	defer file.Close()

	totalSample := 0
	totalDuration := time.Duration(0)
	for {
		samples, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		lenSample := len(samples)
		totalSample += lenSample
		sampleDuration := time.Microsecond * time.Duration(lenSample*125) // 1/8000 = 125 mciro seconds
		totalDuration += sampleDuration

		fmt.Printf("len: %d, time: %f, total time: %f\n", lenSample, sampleDuration.Seconds(), totalDuration.Seconds())
	}
	fmt.Printf("samples: %d\n", totalSample)

	totalLen := 0
	data := make([]byte, 2048)
	for {

		l, err := reader.Read(data)
		if err == io.EOF {
			break
		}

		// lenSample := len(l)
		totalSample += l
		sampleDuration := time.Microsecond * time.Duration(l*125) // 1/8000 = 125 mciro seconds
		totalDuration += sampleDuration

		fmt.Printf("len: %d, time: %f, total time: %f\n", l, sampleDuration.Seconds(), totalDuration.Seconds())

		fmt.Printf("len: %d\n", l)
		totalLen += l
	}

	fmt.Printf("samples: %d\n", totalLen)
}

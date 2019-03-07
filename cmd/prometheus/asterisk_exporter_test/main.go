package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/ivahaev/amigo"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	amiAddress  = flag.String("ami-address", "127.0.0.1", "AMI address")
	amiPort     = flag.String("ami-port", "5038", "AMI port number")
	amiUsername = flag.String("ami-username", "asterisk", "AMI username")
	amiPassword = flag.String("ami-password", "asterisk", "AMI password")

	astNewChannelCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "asterisk_channels",
		Help: "Created new channel count",
	})

	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The toal number of processed events",
	})
)

func initAmiHandle() {
	fmt.Println("Fired initAmigo.")

	fmt.Printf("Connecting to Asterisk AMI. host: %s, port: %s, username: %s, password: -\n", *amiAddress, *amiPort, *amiUsername)
	setting := &amigo.Settings{Username: *amiUsername, Password: *amiPassword, Host: *amiAddress, Port: *amiPort}
	amiHandle := amigo.New(setting)

	amiHandle.Connect()

	// Listen for connection events
	amiHandle.On("connect", func(message string) {
		fmt.Println("Connected.", message)
	})
	amiHandle.On("error", func(message string) {
		fmt.Println("Connection error:", message)
	})

	// registering handler function for event "Newchannel"
	amiHandle.RegisterHandler("Newchannel", func(m map[string]string) {
		astNewChannelCounter.Inc()
	})

	// registering handler function for event "Newchannel"
	amiHandle.RegisterHandler("Newchannel", func(m map[string]string) {
		astNewChannelCounter.Inc()
	})

	amiHandle.RegisterHandler("CoreShowChannel", func(m map[string]string) {
		astNewChannelCounter.Inc()
	})

	amiHandle.RegisterHandler("Hangup", func(m map[string]string) {
		astNewChannelCounter.Desc()
	})

	// check if connected with asterisk, will send action "QueueSummary"
	if amiHandle.Connected() {
		amiHandle.Action(map[string]string{
			"Action":   "QueueSummary",
			"ActionID": "Init",
		})

		amiHandle.Action(map[string]string{
			"Action":   "CoreShowChannels",
			"ActionID": "Init",
		})

	}
}

func amiHandler() {
	fmt.Println("Fired amiHandler")
}

func initProm() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9600", nil)
}

func main() {

	flag.Parse()

	go func() {
		initAmiHandle()
	}()

	go func() {
		initProm()
	}()

	ch := make(chan bool)
	<-ch
}

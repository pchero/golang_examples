package main

import (
	"fmt"

	"github.com/ivahaev/amigo"
)

// Creating handler functions
func DeviceStateChangeHandler(m map[string]string) {
	fmt.Printf("DeviceStateChange event received: %v\n", m)
}

func DefaultHandler(m map[string]string) {
	// fmt.Printf("Event received: %v\n", m)
	fmt.Printf("Event: %s\n", m["Event"])
}

func main() {
	fmt.Println("Init Amigo")

	settings := &amigo.Settings{
		Username: "asterisk",
		Password: "asterisk",
		Host:     "localhost",
	}

	a := amigo.New(settings)

	a.Connect()

	// Listen for connection events
	a.On("connect", func(message string) {
		fmt.Println("Connected", message)
	})

	a.On("error", func(message string) {
		fmt.Println("Connection error:", message)
	})

	// Registering handler function for event "DeviceStateChange"
	a.RegisterHandler("DeviceStateChange", DeviceStateChangeHandler)

	// Registereing default handler function for all events.
	a.RegisterDefaultHandler(DefaultHandler)

	// Optionally create channel to receiving all events
	// and set created channel to receive all events
	c := make(chan map[string]string, 100)
	a.SetEventChannel(c)

	// Check if connected with Asterisk, will send Action "QueueSummary"
	if a.Connected() {
		result, err := a.Action(map[string]string{"Action": "QueueSummary", "ActionID": "Init"})
		// If not error, processing result. Response on Action will follow in defined events.
		// You need to catch them in event channel, DefaultHandler or specified HandlerFunction
		fmt.Println(result, err)
	}

	ch := make(chan bool)
	<-ch
}

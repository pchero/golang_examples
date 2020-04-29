package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const input = `
{
	"type": "sound",
	"msg": {
		"description": "dynamite",
		"authority": "the Bruce Dickinson"
	}
}
`

type Envelope struct {
	Type string
	Msg  interface{}
}

type Sound struct {
	Description string
	Authority   string
}

func main() {
	var msg json.RawMessage
	env := Envelope{
		Msg: &msg,
	}
	if err := json.Unmarshal([]byte(input), &env); err != nil {
		log.Fatal(err)
	}
	switch env.Type {
	case "sound":
		var s Sound
		if err := json.Unmarshal(msg, &s); err != nil {
			log.Fatal(err)
		}
		var desc string = s.Description
		fmt.Println(desc)
	default:
		log.Fatalf("unknown message type: %q", env.Type)
	}
}

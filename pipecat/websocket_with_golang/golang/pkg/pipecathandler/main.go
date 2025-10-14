package pipecathandler

import (
	"net"
	"net/http"
	"os/exec"
)

type pipecatHandler struct {
	llm         string
	stt         string
	tts         string
	messageFile string

	// pipecat
	listener net.Listener
	port     int
	server   *http.Server
	cmd      *exec.Cmd // python

	// asterisk
}

type PipecatcatHandler interface {
	Start(llm string, stt string, tts string, message_file string) error
	Stop()

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func NewPipecatHandler() PipecatcatHandler {
	return &pipecatHandler{}
}

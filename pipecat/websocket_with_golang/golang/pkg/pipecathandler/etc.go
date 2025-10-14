package pipecathandler

import (
	"net"
	"net/http"
)

func (h *pipecatHandler) setPort(port int) {
	h.port = port
}

func (h *pipecatHandler) setServer(server *http.Server) {
	h.server = server
}

func (h *pipecatHandler) setListener(listener net.Listener) {
	h.listener = listener
}

func (h *pipecatHandler) setArgs(llm string, stt string, tts string, message_file string) {
	h.llm = llm
	h.stt = stt
	h.tts = tts
	h.messageFile = message_file
}

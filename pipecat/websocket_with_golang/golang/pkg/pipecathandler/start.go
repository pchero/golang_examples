package pipecathandler

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (h *pipecatHandler) Start(llm string, stt string, tts string, message_file string) error {
	log := logrus.WithFields(logrus.Fields{
		"func": "Start",
	})
	log.Debugf("Starting Pipecat WebSocket server...")

	// set args for debug
	h.setArgs(llm, stt, tts, message_file)

	if errWebsocket := h.startWebsocketListener(); errWebsocket != nil {
		log.Errorf("Error starting Pipecat WebSocket server: %v", errWebsocket)
		return errors.Wrapf(errWebsocket, "could not start the pipecat websocket server")
	}

	if errPython := h.startPythonRunner(llm, stt, tts, message_file); errPython != nil {
		log.Errorf("Error starting Pipecat Python runner: %v", errPython)
		return errors.Wrapf(errPython, "could not start the pipecat python runner")
	}

	return nil
}

func (h *pipecatHandler) Stop() {
	log := logrus.WithFields(logrus.Fields{
		"func": "Stop",
	})

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if h.server != nil {
		log.Debugf("Shutting down Go WebSocket server...")
		if errShutdown := h.server.Shutdown(shutdownCtx); errShutdown != nil {
			log.Errorf("Go WebSocket server shutdown failed: %v", errShutdown)
		}
		log.Debugf("Go WebSocket server shut down.")
	}

	if h.cmd != nil && h.cmd.Process != nil {
		log.Debugf("Terminating Python process with PID %d...", h.cmd.Process.Pid)
		if errKill := h.cmd.Process.Kill(); errKill != nil {
			log.Errorf("Failed to forcefully kill Python process: %v", errKill)
		}
		log.Debugf("Python process terminated.")
	}

	h.listener.Close()
}

func (h *pipecatHandler) startWebsocketListener() error {
	log := logrus.WithFields(logrus.Fields{
		"func": "startPipecatWebsocket",
	})

	app := http.NewServeMux()
	app.Handle("/ws", h)

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Errorf("Failed to listen on ephemeral port: %v", err)
		return errors.Wrapf(err, "failed to listen on ephemeral port")
	}
	h.setListener(listener)

	port := listener.Addr().(*net.TCPAddr).Port
	h.setPort(port)
	log.Debugf("Server assigned to port: %d", port)

	server := &http.Server{
		Handler: app,
	}
	h.setServer(server)

	go func() {
		log.Debugf("Starting HTTP server on %s", listener.Addr().String())
		if errServe := server.Serve(listener); errServe != nil && errServe != http.ErrServerClosed {
			log.Errorf("Could not start HTTP server: %v", errServe)
		}
		log.Debugf("HTTP server stopped")
	}()

	return nil
}

func (h *pipecatHandler) startPythonRunner(llm string, stt string, tts string, message_file string) error {
	log := logrus.WithFields(logrus.Fields{
		"func": "startPipecatPython",
	})

	pythonInterpreter := "scripts/pipecat/venv/bin/python3"
	pythonScript := "scripts/pipecat/main.py"

	args := []string{
		"--ws_server_url", fmt.Sprintf("ws://localhost:%d/ws", h.port),
		"--llm", llm,
		"--stt", stt,
		"--tts", tts,
		"--messages_file", message_file,
	}

	cmdArgs := append([]string{pythonScript}, args...)
	cmd := exec.Command(pythonInterpreter, cmdArgs...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "could to get stdout pipe for python")
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrapf(err, "could to get stderr pipe for python")
	}

	log.Debugf("Executing pipecat python script: %s %v", pythonInterpreter, cmdArgs)
	if errStart := cmd.Start(); errStart != nil {
		return errors.Wrapf(errStart, "could not start python client")
	}
	log.Debugf("Python client process started with PID: %d", cmd.Process.Pid)

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Debugf("[PYTHON-CLIENT-STDOUT] %s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Errorf("Error reading Python client stdout: %v", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Debugf("[PYTHON-CLIENT-STDERR] %s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Errorf("Error reading Python client stderr: %v", err)
		}
	}()

	go func() {
		// wait for the python process to exit
		if errPython := cmd.Wait(); errPython != nil {
			log.Errorf("Python client process exited with error: %v", errPython)
		}
	}()

	return nil
}

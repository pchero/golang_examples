package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 파이썬 스크립트 경로 및 인자 설정
	pythonScript := "../pipecat/server.py"
	args := []string{
		pythonScript,
		"--messages_file", "./messages.json",
		"--llm", "openai.gpt-4o-mini",
		"--stt", "deepgram",
		"--tts", "cartesia",
	}

	// 2. 파이썬 인터프리터 경로 설정
	// 가상 환경의 파이썬 인터프리터를 사용하는 것이 가장 좋습니다.
	// 님의 경우 venv/bin/python3
	pythonInterpreter := "./venv/bin/python3" // 프로젝트 루트 기준 상대 경로

	// exec.Command를 사용하여 명령 생성
	cmd := exec.Command(pythonInterpreter, args...)

	// 파이썬 스크립트가 실행될 디렉토리 설정 (선택 사항, 하지만 일반적으로 필요)
	// 님의 경우 server.py가 있는 'pipecat' 디렉토리가 될 것입니다.
	// 현재 Go 프로그램이 pipecat 디렉토리에서 실행된다면 필요 없습니다.
	// 그러나 Go 프로그램이 다른 디렉토리에서 실행될 경우 정확한 경로를 지정해야 합니다.
	// cmd.Dir = "/home/pchero/github/golang_examples/pipecat/websocket_example/pipecat"

	log.Printf("Starting Python server: %s %v", pythonInterpreter, args)

	// 3. 표준 출력(stdout) 및 표준 에러(stderr) 스트림 캡처
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Failed to get stderr pipe: %v", err)
	}

	// 4. 파이썬 스크립트 시작
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Python server: %v", err)
	}
	log.Printf("Python server started with PID: %d", cmd.Process.Pid)

	// 5. stdout과 stderr을 읽기 위한 고루틴 시작
	// 파이썬 서버의 출력을 실시간으로 Go 프로그램 콘솔에 표시
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Printf("[PYTHON-STDOUT] %s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading stdout: %v", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Printf("[PYTHON-STDERR] %s\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading stderr: %v", err)
		}
	}()

	// 6. Ctrl+C (SIGINT) 신호를 받아서 자식 프로세스를 종료하도록 처리
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Interrupt signal received. Terminating Python server...")
		// 자식 프로세스에 SIGTERM 신호를 보내서 종료를 요청합니다.
		// 파이썬 스크립트가 SIGTERM을 제대로 처리한다면 깔끔하게 종료됩니다.
		// 아니면 Kill()을 사용할 수도 있지만, Kill()은 강제 종료입니다.
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Error sending SIGTERM to Python process: %v", err)
			// SIGTERM 실패 시 강제 종료 (fallback)
			if err := cmd.Process.Kill(); err != nil {
				log.Fatalf("Failed to kill Python process: %v", err)
			}
		}

		// 자식 프로세스가 종료될 때까지 기다림
		// 일정 시간 (예: 5초) 기다린 후에도 종료되지 않으면 강제 종료
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil {
				log.Printf("Python server exited with error: %v", err)
			} else {
				log.Println("Python server terminated gracefully.")
			}
		case <-time.After(5 * time.Second): // 5초 타임아웃
			log.Println("Python server did not terminate in time, killing it forcefully.")
			if err := cmd.Process.Kill(); err != nil {
				log.Fatalf("Failed to kill Python process: %v", err)
			}
			<-done // Kill 후에도 Wait()를 호출하여 리소스 정리
		}

	case <-time.After(30 * time.Second): // 예시: 30초 후 자동으로 파이썬 서버 종료 (옵션)
		log.Println("Go program timeout. Terminating Python server...")
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Error sending SIGTERM to Python process: %v", err)
			if err := cmd.Process.Kill(); err != nil {
				log.Fatalf("Failed to kill Python process: %v", err)
			}
		}
		cmd.Wait() // 종료될 때까지 기다림
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			timestamp := time.Now().Format("2006-01-02 15:04:05")

			logEntry := fmt.Sprintf("[%s] INPUT: %s", timestamp, input)
			logger.Println(logEntry)

			if strings.ToLower(strings.TrimSpace(input)) == "quit" ||
				strings.ToLower(strings.TrimSpace(input)) == "exit" {
				logger.Println("Quit command received, exiting...")
				done <- true
				return
			}
		}

		if scanner.Err() != nil {
			logger.Printf("Error reading input: %v", scanner.Err())
		}
	}()

	go func() {
		sig := <-sigs
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logger.Printf("[%s] Received signal: %v", timestamp, sig)
		done <- true
	}()

	logger.Println("Application started. Type messages (they will be saved to app.log)")
	logger.Println("Type 'quit' or 'exit' to stop, or press Ctrl+C")

	<-done

	logger.Println("Application exiting")
}

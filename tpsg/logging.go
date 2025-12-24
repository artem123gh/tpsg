package tpsg

import (
	"fmt"
	"time"
)

func LogInfo(message string) {
	timestamp := time.Now().Format("2006.01.02 15:04:05.000")
	fmt.Printf("Info  | %s | %s\n", timestamp, message)
}

func LogEvent(message string) {
	timestamp := time.Now().Format("2006.01.02 15:04:05.000")
	fmt.Printf("Event | %s | %s\n", timestamp, message)
}

func LogError(message string) {
	timestamp := time.Now().Format("2006.01.02 15:04:05.000")
	fmt.Printf("Error | %s | %s\n", timestamp, message)
}

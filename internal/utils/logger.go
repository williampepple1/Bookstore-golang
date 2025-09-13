package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// LogLevel represents the log level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var currentLogLevel LogLevel = INFO

// SetLogLevel sets the current log level
func SetLogLevel(level LogLevel) {
	currentLogLevel = level
}

// LogInfo logs an info message
func LogInfo(message string, data interface{}) {
	if currentLogLevel <= INFO {
		logMessage("INFO", message, data)
	}
}

// LogError logs an error message
func LogError(message string, data interface{}) {
	if currentLogLevel <= ERROR {
		logMessage("ERROR", message, data)
	}
}

// LogWarn logs a warning message
func LogWarn(message string, data interface{}) {
	if currentLogLevel <= WARN {
		logMessage("WARN", message, data)
	}
}

// LogDebug logs a debug message
func LogDebug(message string, data interface{}) {
	if currentLogLevel <= DEBUG {
		logMessage("DEBUG", message, data)
	}
}

// logMessage logs a message with the specified level
func logMessage(level, message string, data interface{}) {
	timestamp := getCurrentTimestamp()
	
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("[%s] %s %s - %v", timestamp, level, message, data)
		} else {
			log.Printf("[%s] %s %s - %s", timestamp, level, message, string(jsonData))
		}
	} else {
		log.Printf("[%s] %s %s", timestamp, level, message)
	}
}

// getCurrentTimestamp returns the current timestamp
func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", getCurrentTime().Unix())
}

// getCurrentTime returns the current time (can be mocked for testing)
func getCurrentTime() time.Time {
	return time.Now()
}

// LogRequest logs HTTP request details
func LogRequest(c interface{}, duration time.Duration, err error) {
	// This is a simplified version for now
	// In a real implementation, you'd extract the details from the fiber context
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"duration":  duration.String(),
		"error":     err != nil,
	}
	
	// Log based on error
	if err != nil {
		LogError("HTTP Request", logEntry)
	} else {
		LogInfo("HTTP Request", logEntry)
	}
}

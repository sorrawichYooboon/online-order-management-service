package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var sensitiveKeys = []string{
	"password", "token", "apikey", "secret", "accessToken", "refreshToken",
}

func InitLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(os.Stdout)
	LogInfo("InitLogger", "Logger initialized")
}

func Log(severity string, sectionFunc string, message string, fields ...map[string]any) {
	var sanitizedFields map[string]any
	if len(fields) > 0 {
		sanitizedFields = sanitizeFields(fields[0])
	} else {
		sanitizedFields = map[string]any{}
	}

	logMessage := fmt.Sprintf("[%s] [%s] %s | data: %+v", strings.ToUpper(severity), sectionFunc, message, sanitizedFields)
	log.Println(logMessage)
}

func LogInfo(sectionFunc string, message string, fields ...map[string]any) {
	if len(fields) > 0 {
		Log("INFO", sectionFunc, message, fields[0])
	} else {
		Log("INFO", sectionFunc, message)
	}
}

func LogError(sectionFunc string, err error, fields ...map[string]any) {
	if err != nil {
		logFields := map[string]any{"error": err.Error()}
		if len(fields) > 0 {
			for k, v := range fields[0] {
				logFields[k] = v
			}
		}
		Log("ERROR", sectionFunc, err.Error(), logFields)
	}
}

func sanitizeFields(fields map[string]any) map[string]any {
	sanitized := make(map[string]any)
	for key, value := range fields {
		if isSensitiveKey(key) {
			sanitized[key] = "*****"
		} else {
			sanitized[key] = value
		}
	}
	return sanitized
}

func isSensitiveKey(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(lowerKey, sensitive) {
			return true
		}
	}
	return false
}

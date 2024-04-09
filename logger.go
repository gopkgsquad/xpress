package xpress

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the log level.
type LogLevel int

const (
	// LogLevelDebug represents debug log level.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents info log level.
	LogLevelInfo
	// LogLevelWarning represents warning log level.
	LogLevelWarning
	// LogLevelError represents error log level.
	LogLevelError
	// LogLevelFatal represents fatal log level.
	LogLevelFatal
)

type LogFormat int

const (
	FormatColorized LogFormat = iota
	FormatJSON
)

// Logger represents a custom logger with different log levels.
type Logger struct {
	level            LogLevel
	logger           *log.Logger
	timeColor        string
	fileColor        string
	renderCallerInfo *bool
	format           LogFormat
}

func (l *Logger) LogHTTPRequest(level LogLevel, message interface{}) {
	if l.level <= level {
		var logMessage string
		switch l.format {
		case FormatColorized:
			timestamp := fmt.Sprintf("%s[%s]%s ", l.timeColor, time.Now().Format("2006/01/02 15:04:05"), colorReset)
			levelColor := l.getLevelColor(level)
			logMessage = timestamp + levelColor + fmt.Sprint(message) + colorReset

		case FormatJSON:
			logStruct := map[string]interface{}{
				"timestamp": time.Now().Format("2006/01/02 15:04:05"),
				"level":     l.getLevelString(level),
			}
			messageMap, ok := message.(map[string]interface{})
			if ok {
				for k, v := range messageMap {
					logStruct[k] = v
				}
			} else {
				logStruct["message"] = message
			}
			if l.renderCallerInfo != nil && *l.renderCallerInfo {
				logStruct["caller"] = l.getCallerInfoPlain()
			}
			bytes, err := json.Marshal(logStruct)
			if err != nil {
				logMessage = fmt.Sprintf("Failed to marshal log message to JSON: %v", err)
			} else {
				logMessage = string(bytes)
			}
		}
		l.logger.Output(3, logMessage)
	}
}

func (l *Logger) Log(level LogLevel, message interface{}) {
	if l.level <= level {
		var logMessage string
		switch l.format {
		case FormatColorized:
			timestamp := fmt.Sprintf("%s[%s]%s ", l.timeColor, time.Now().Format("2006/01/02 15:04:05"), colorReset)
			levelColor := l.getLevelColor(level)
			logMessage = timestamp + levelColor + fmt.Sprint(message) + colorReset
			if l.renderCallerInfo != nil && *l.renderCallerInfo {
				caller := l.getCallerInfo()
				callerInfo := fmt.Sprintf("%s%s%s ", l.fileColor, caller, colorReset)
				logMessage = callerInfo + logMessage
			}
		case FormatJSON:
			logStruct := struct {
				Timestamp string      `json:"timestamp"`
				Level     string      `json:"level"`
				Message   interface{} `json:"message"`
				Caller    string      `json:"caller,omitempty"`
			}{
				Timestamp: time.Now().Format("2006/01/02 15:04:05"),
				Level:     l.getLevelString(level),
				Message:   message,
			}
			if l.renderCallerInfo != nil && *l.renderCallerInfo {
				logStruct.Caller = l.getCallerInfoPlain()
			}
			bytes, err := json.Marshal(logStruct)
			if err != nil {
				logMessage = fmt.Sprintf("Failed to marshal log message to JSON: %v", err)
			} else {
				logMessage = string(bytes)
			}
		}
		l.logger.Output(3, logMessage)
	}
}

func (l *Logger) getCallerInfoPlain() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}
	return fmt.Sprintf("[%s:%d]", file, line)
}

func (l *Logger) getLevelString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg interface{}) {
	l.Log(LogLevelDebug, msg)
}

// Info logs an info message.
func (l *Logger) Info(msg interface{}) {
	l.Log(LogLevelInfo, msg)
}

// Warning logs a warning message.
func (l *Logger) Warning(msg interface{}) {
	l.Log(LogLevelWarning, msg)
}

// Error logs an error message.
func (l *Logger) Error(msg interface{}) {
	l.Log(LogLevelError, msg)
}

// Fatal logs a fatal message and exits the application.
func (l *Logger) Fatal(msg interface{}) {
	l.Log(LogLevelFatal, msg)
	os.Exit(1)
}

// Infof logs an info message with formatting.
func (l *Logger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Log(LogLevelInfo, message)
}

// Warningf logs a warning message with formatting.
func (l *Logger) Warningf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Log(LogLevelWarning, message)
}

// Errorf logs an error message with formatting.
func (l *Logger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Log(LogLevelError, message)
}

// Fatalf logs a fatal message with formatting and exits the application.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Log(LogLevelFatal, message)
	os.Exit(1)
}

func (l *Logger) getCallerInfo() string {
	_, file, line, _ := runtime.Caller(3) // Adjusted depth to get caller outside logger functions
	file = filepath.Base(file)
	return fmt.Sprintf("%s[%s:%d]%s", l.fileColor, file, line, colorReset)
}

func (l *Logger) getLevelColor(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return colorBlue
	case LogLevelInfo:
		return colorGreen
	case LogLevelWarning:
		return colorYellow
	case LogLevelError, LogLevelFatal:
		return colorRed
	default:
		return colorReset
	}
}

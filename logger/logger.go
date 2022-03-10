package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	// Init the logger
	Logger = logrus.New()

	// Set formatting
	Logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	// Set the log level
	Logger.SetLevel(logrus.InfoLevel)

	// Set the output of the logger
	Logger.Out = os.Stdout
}

// InfoMessage displays info message
func InfoMessage(f string, a ...interface{}) {
	message := fmt.Sprintf(f, a...)
	Logger.Info(message)
}

// WarnMessage displays warning message
func WarnMessage(f string, a ...interface{}) {
	message := fmt.Sprintf(f, a...)
	Logger.Warn(message)

}

// ErrorMessage displays error message
func ErrorMessage(f string, a ...interface{}) {
	message := fmt.Sprintf(f, a...)
	_, file, line, _ := runtime.Caller(1)
	Logger.WithField("file", fmt.Sprintf("%s %d", file, line)).Error(message)

}

// FatalMessage displays fatal messages
func FatalMessage(f string, a ...interface{}) {
	message := fmt.Sprintf(f, a...)

	_, file, line, _ := runtime.Caller(1)
	Logger.WithField("file", fmt.Sprintf("%s %d", file, line)).Fatal(message)

}

// DebugMessage displays debugging messages
func DebugMessage(f string, a ...interface{}) {
	message := fmt.Sprintf(f, a...)

	_, file, line, _ := runtime.Caller(1)
	Logger.WithField("file", fmt.Sprintf("%s %d", file, line)).Debug(message)
}

// TraceMessage displays trace messages
func TraceMessage(f string, a ...interface{}) {
	message := fmt.Sprintf(f, a...)

	_, file, line, _ := runtime.Caller(1)
	Logger.WithField("file", fmt.Sprintf("%s %d", file, line)).Trace(message)
}

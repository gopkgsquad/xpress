package xpress

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	log.SetReportCaller(true)

	timestamp := "2nd Jan,2006 03:04 PM"

	// Extract the day suffix
	daySuffix := ""
	day := strings.Split(timestamp, " ")[0]
	dayNum, _ := strconv.Atoi(day)
	switch dayNum {
	case 1, 21, 31:
		daySuffix = "st"
	case 2, 22:
		daySuffix = "nd"
	case 3, 23:
		daySuffix = "rd"
	default:
		daySuffix = "th"
	}

	// Replace "th" with "th," "st" with "st," etc.
	day = day[:len(day)-2] + daySuffix

	// Reconstruct the timestamp
	timestamp = day + " " + strings.Join(strings.Split(timestamp, " ")[1:], " ")

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		TimestampFormat:  timestamp,
		FullTimestamp:    true,
		CallerPrettyfier: caller(),
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFile: "caller",
		},
	})

	return log
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {

		return "", ""
	}
}

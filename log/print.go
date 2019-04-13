package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

//Print print something
func Print(level string, msg ...interface{}) {
	// 行數位置
	_, fn, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("Log error: %+v, %+v, %+v, \n", fn, line, ok)
		return
	}

	// Add to field
	file := strings.Split(fn, "/")
	logField := logrus.Fields{
		"file": file[len(file)-2] + "/" + file[len(file)-1] + " " + strconv.Itoa(line),
	}

	switch level {
	case "info":
		logrus.Info(msg)
	case "warn":
		logrus.Warn(msg)
	case "error":
		logrus.WithFields(logField).Error(msg)
	default:
		logrus.WithFields(logField).Error(msg)
	}
}

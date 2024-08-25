package logger

import (
	"log"
	"os"
	"strings"
)

var (
	debugLogger *log.Logger
	errorLogger *log.Logger
	infoLogger  *log.Logger
)

func init() {
	debugLogger = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
}

func Info(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	d, _ := os.LookupEnv("DEBUG")
	if strings.EqualFold(d, "true") {
		debugLogger.Printf(format, v...)
	}
}

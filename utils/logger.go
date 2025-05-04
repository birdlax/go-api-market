package utils

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

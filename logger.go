package Grawler

// Wrapper of logger functions

import (
	"log"
)

func logInfo(str string) {
	log.Printf("[INFO]: %v", str)
}

func logFatal(str string) {
	log.Fatalf("[ERROR]: %v", str)
}

package Grawler

// Wrapper of log functions

import (
	"log"
)

func logInfo(str string) {
	log.Printf("[INFO]: %v", str)
}

func logFatal(str string) {
	log.Fatalf("[ERROR]: %v", str)
}

package services

import (
	"log"
)

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

package core

import (
	"fmt"
	"log"
)

func AssertErrorToNilf(message string, err error) bool {
	if err != nil {
		log.Fatalf(message, err)
		return true
	}
	return false
}

func LogLine(message string, args ...any) {
	output := fmt.Sprintf(message, args...)
	fmt.Printf("%s \n", output)
}

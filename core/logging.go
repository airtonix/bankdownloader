package core

import "log"

func AssertErrorToNilf(message string, err error) bool {
	if err != nil {
		log.Fatalf(message, err)
		return true
	}
	return false
}

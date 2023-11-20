package core

import "strings"

func UnQuote(str string) string {
	if len(str) < 2 {
		return str
	}

	if str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}

// returns a string of stars instead of the characters in the string
func Stars(value string) string {
	return strings.Repeat("*", len(value))
}

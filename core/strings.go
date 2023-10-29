package core

func UnQuote(str string) string {
	if len(str) < 2 {
		return str
	}

	if str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}

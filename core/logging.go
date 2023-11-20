package core

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
)

func AssertErrorToNilf(message string, err error) bool {
	if err != nil {
		logrus.Panic(
			color.FgRed.Render(fmt.Sprintf(message, err)),
		)
		return true
	}
	return false
}

func Header(message string) {
	fmt.Printf("\n\n%s\n\n",
		color.FgCyan.Render(message),
	)
}

func KeyValue(key string, value any) {
	fmt.Printf(
		"\t%s: %s\n",
		color.FgYellow.Render(key),
		color.FgWhite.Render(value),
	)
}

func Action(name string) {
	fmt.Printf(
		"\t%s",
		color.FgGray.Render(name),
	)
}

package utils

import (
	"strings"
)

func CountLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

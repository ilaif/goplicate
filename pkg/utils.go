package pkg

import (
	"strings"
)

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

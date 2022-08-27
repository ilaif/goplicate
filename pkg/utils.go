package pkg

import (
	"strings"
)

func splitLines(bytes []byte) []string {
	return strings.Split(string(bytes), "\n")
}

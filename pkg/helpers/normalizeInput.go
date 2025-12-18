package helpers

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func NormalizeInput(s string) string {
	// Trim spaces & lowercase
	s = strings.TrimSpace(strings.ToLower(s))

	// Normalize Unicode (remove accents)
	t := norm.NFD.String(s)
	sb := strings.Builder{}
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue // remove accent marks
		}
		sb.WriteRune(r)
	}

	return sb.String()
}

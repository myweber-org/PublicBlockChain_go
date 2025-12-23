package csvutils

import (
	"strings"
)

func SanitizeCSVRow(row []string) []string {
	sanitized := make([]string, len(row))
	for i, field := range row {
		sanitized[i] = strings.TrimSpace(field)
	}
	return sanitized
}
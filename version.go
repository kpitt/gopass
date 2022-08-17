package main

import (
	"fmt"
	"strings"
)

func FormatVersion(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}

	return version + dateStr
}

package util

import (
	"regexp"
	"strings"
)

var (
	illegalChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

	// Reserved names on Windows
	reservedNames = map[string]bool{
		"CON": true, "PRN": true, "AUX": true, "NUL": true,
		"COM1": true, "COM2": true, "COM3": true, "COM4": true,
		"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true,
	}
)

func SanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")

	name = illegalChars.ReplaceAllString(name, "_")

	name = strings.TrimRight(name, ". ")

	// Prevent empty name
	if name == "" {
		name = "file"
	}

	if reservedNames[strings.ToUpper(name)] {
		name = "_" + name
	}

	if len(name) > 255 {
		name = name[:255]
	}

	return name
}

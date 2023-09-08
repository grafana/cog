package tools

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func UpperCamelCase(s string) string {
	s = LowerCamelCase(s)

	// Uppercase the first letter
	if len(s) > 0 {
		s = strings.ToUpper(s[:1]) + s[1:]
	}

	return s
}

func LowerCamelCase(s string) string {
	// Replace all underscores/dashes with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	// Title case s
	s = cases.Title(language.AmericanEnglish, cases.NoLower).String(s)

	// Remove all spaces
	s = strings.ReplaceAll(s, " ", "")

	// Lowercase the first letter
	if len(s) > 0 {
		s = strings.ToLower(s[:1]) + s[1:]
	}

	return s
}

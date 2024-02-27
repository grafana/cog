package tools

import (
	"regexp"
	"strings"

	"github.com/huandu/xstrings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var nonAlphaNumRegex = regexp.MustCompile("[^a-zA-Z0-9 ]+")

func UpperSnakeCase(s string) string {
	return strings.ToUpper(xstrings.ToSnakeCase(s))
}

func SnakeCase(s string) string {
	return xstrings.ToSnakeCase(s)
}

func UpperCamelCase(s string) string {
	s = LowerCamelCase(s)

	// Uppercase the first letter
	if len(s) > 0 {
		s = strings.ToUpper(s[:1]) + s[1:]
	}

	return s
}

func LowerCamelCase(s string) string {
	// Replace all non-alphanumeric characters by spaces
	s = nonAlphaNumRegex.ReplaceAllString(s, " ")

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

// CleanupNames removes all non-alphanumeric characters
func CleanupNames(s string) string {
	return nonAlphaNumRegex.ReplaceAllString(s, "")
}

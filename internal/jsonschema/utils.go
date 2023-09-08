package jsonschema

import (
	"strings"

	schemaparser "github.com/santhosh-tekuri/jsonschema"
)

func stringInList(list []string, input string) bool {
	for _, value := range list {
		if value == input {
			return true
		}
	}

	return false
}

func schemaComments(schema *schemaparser.Schema) []string {
	comment := schema.Description

	lines := strings.Split(comment, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered
}

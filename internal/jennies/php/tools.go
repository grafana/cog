package php

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grafana/cog/internal/tools"
)

func formatPackageName(pkg string) string {
	rgx := regexp.MustCompile("[^a-zA-Z0-9_]+")

	return tools.UpperCamelCase(rgx.ReplaceAllString(pkg, ""))
}

func formatObjectName(name string) string {
	return tools.UpperCamelCase(name)
}

func formatConstantName(name string) string {
	return tools.UpperSnakeCase(name)
}

func formatFieldName(name string) string {
	return tools.LowerCamelCase(name)
}

func formatEnumMemberName(name string) string {
	return tools.LowerCamelCase(name)
}

func formatCommentsBlock(comments []string) string {
	if len(comments) == 0 {
		return ""
	}

	var buffer strings.Builder

	if len(comments) != 0 {
		buffer.WriteString("/**\n")
	}
	for _, commentLine := range comments {
		buffer.WriteString(fmt.Sprintf(" * %s\n", commentLine))
	}
	if len(comments) != 0 {
		buffer.WriteString(" */\n")
	}

	return buffer.String()
}

package python

import (
	"fmt"
	"strings"
)

type raw string

func formatValue(val any) string {
	if val == nil {
		return "None"
	}

	if rawVal, ok := val.(raw); ok {
		return string(rawVal)
	}

	if asBool, ok := val.(bool); ok {
		if asBool {
			return "True"
		}

		return "False"
	}

	if list, ok := val.([]any); ok {
		items := make([]string, 0, len(list))

		for _, item := range list {
			items = append(items, formatValue(item))
		}

		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	}

	return fmt.Sprintf("%#v", val)
}

func escapeFieldName(name string) string {
	if isReservedPythonKeyword(name) {
		return name + "_val"
	}

	return name
}

func isReservedPythonKeyword(input string) bool {
	// see: https://docs.python.org/3/reference/lexical_analysis.html#keywords
	return input == "False" ||
		input == "await" ||
		input == "else" ||
		input == "import" ||
		input == "pass" ||
		input == "None" ||
		input == "break" ||
		input == "except" ||
		input == "in" ||
		input == "raise" ||
		input == "True" ||
		input == "class" ||
		input == "finally" ||
		input == "is" ||
		input == "return" ||
		input == "and" ||
		input == "continue" ||
		input == "for" ||
		input == "lambda" ||
		input == "try" ||
		input == "as" ||
		input == "def" ||
		input == "from" ||
		input == "nonlocal" ||
		input == "while" ||
		input == "assert" ||
		input == "del" ||
		input == "global" ||
		input == "not" ||
		input == "with" ||
		input == "async" ||
		input == "elif" ||
		input == "if" ||
		input == "or" ||
		input == "yield"
}

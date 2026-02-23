package option

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/veneers"
)

type Rule struct {
	Selector *Selector
	Action   *Action
}

func (rule Rule) Matches(builder ast.Builder, option ast.Option) bool {
	return rule.Selector.Matches(builder, option)
}

func (rule Rule) Apply(schemas ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option {
	return rule.Action.run(schemas, builder, option)
}

func (rule Rule) String() string {
	return fmt.Sprintf("selector=%s, action=%s", rule.Selector, rule.Action)
}

type Action struct {
	description string
	run         ActionRunner
}

func (action Action) String() string {
	return action.description
}

func Rename(selector *Selector, newName string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("rename[as: '%s']", newName),
			run:         RenameAction(newName),
		},
	}
}

func ArrayToAppend(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "array_to_append",
			run:         ArrayToAppendAction(),
		},
	}
}

func MapToIndex(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "map_to_index",
			run:         MapToIndexAction(),
		},
	}
}

func RenameArguments(selector *Selector, newNames []string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("rename_arguments[as: (%s)]", strings.Join(newNames, ", ")),
			run:         RenameArgumentsAction(newNames),
		},
	}
}

func Omit(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "omit",
			run:         OmitAction(),
		},
	}
}

func VeneerTrailAsComments(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "veneer_trail_as_comments",
			run:         VeneerTrailAsCommentsAction(),
		},
	}
}

func UnfoldBoolean(selector *Selector, unfoldOpts BooleanUnfold) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "unfold_boolean",
			run:         UnfoldBooleanAction(unfoldOpts),
		},
	}
}

func StructFieldsAsArguments(selector *Selector, explicitFields ...string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "struct_fields_as_arguments",
			run:         StructFieldsAsArgumentsAction(explicitFields...),
		},
	}
}

func StructFieldsAsOptions(selector *Selector, explicitFields ...string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "struct_fields_as_options",
			run:         StructFieldsAsOptionsAction(explicitFields...),
		},
	}
}

func DisjunctionAsOptions(selector *Selector, argumentIndex int) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "disjunction_as_options",
			run:         DisjunctionAsOptionsAction(argumentIndex),
		},
	}
}

func Duplicate(selector *Selector, duplicateName string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("duplicate[as: '%s']", duplicateName),
			run:         DuplicateAction(duplicateName),
		},
	}
}

func AddAssignment(selector *Selector, assignment veneers.Assignment) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "add_assignment",
			run:         AddAssignmentAction(assignment),
		},
	}
}

func AddComments(selector *Selector, comments []string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "add_comments",
			run:         AddCommentsAction(comments),
		},
	}
}

func Debug(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "debug",
			run:         DebugAction(),
		},
	}
}

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

// Rename renames an option.
func Rename(selector *Selector, newName string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("rename[as: '%s']", newName),
			run:         RenameAction(newName),
		},
	}
}

// ArrayToAppend updates the option to perform an "append" assignment.
//
// Example:
//
//	```
//	func Tags(tags []string) {
//		this.resource.tags = tags
//	}
//	```
//
// Will become:
//
//	```
//	func Tags(tags string) {
//		this.resource.tags.append(tags)
//	}
//	```
//
// This action returns the option unchanged if:
//   - it doesn't have exactly one argument
//   - the argument is not an array
func ArrayToAppend(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "array_to_append",
			run:         ArrayToAppendAction(),
		},
	}
}

// MapToIndex updates the option to perform an "index" assignment.
//
// Example:
//
//	```
//	func Elements(elements map[string]Element) {
//		this.resource.elements = elements
//	}
//	```
//
// Will become:
//
//	```
//	func Elements(key string, elements Element) {
//		this.resource.elements[key] = tags
//	}
//	```
//
// This action returns the option unchanged if:
//   - it doesn't have exactly one argument
//   - the argument is not a map
func MapToIndex(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "map_to_index",
			run:         MapToIndexAction(),
		},
	}
}

// RenameArguments renames the arguments of an options.
func RenameArguments(selector *Selector, newNames []string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("rename_arguments[as: (%s)]", strings.Join(newNames, ", ")),
			run:         RenameArgumentsAction(newNames),
		},
	}
}

// Omit removes an option.
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

// UnfoldBoolean transforms an option accepting a boolean argument into two argument-less options.
//
// Example:
//
//	```
//	func Editable(editable bool) {
//		this.resource.editable = editable
//	}
//	```
//
// Will become:
//
//	```
//	func Editable() {
//		this.resource.editable = true
//	}
//
//	func ReadOnly() {
//		this.resource.editable = false
//	}
//	```
func UnfoldBoolean(selector *Selector, unfoldOpts BooleanUnfold) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "unfold_boolean",
			run:         UnfoldBooleanAction(unfoldOpts),
		},
	}
}

// StructFieldsAsArguments uses the fields of the first argument's struct (assuming it is one) and turns them
// into arguments.
//
// Optionally, an explicit list of fields to turn into arguments can be given.
//
// Example:
//
//	```
//	func Time(time {from string, to string) {
//		this.resource.time = time
//	}
//	```
//
// Will become:
//
//	```
//	func Time(from string, to string) {
//		this.resource.time.from = from
//		this.resource.time.to = to
//	}
//	```
//
// This action returns the option unchanged if:
//   - it has no arguments
//   - the first argument is not a struct or a reference to one
//
// FIXME: considers the first argument only.
func StructFieldsAsArguments(selector *Selector, explicitFields ...string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "struct_fields_as_arguments",
			run:         StructFieldsAsArgumentsAction(explicitFields...),
		},
	}
}

// StructFieldsAsOptions uses the fields of the first argument's struct (assuming it is one) and turns them
// into options.
//
// Optionally, an explicit list of fields to turn into options can be given.
//
// Example:
//
//	```
//	func GridPos(gridPos {x int, y int) {
//		this.resource.gridPos = gridPos
//	}
//	```
//
// Will become:
//
//	```
//	func X(x int) {
//		this.resource.gridPos.x = x
//	}
//
//	func Y(y int) {
//		this.resource.gridPos.y = y
//	}
//	```
//
// This action returns the option unchanged if:
//   - it has no arguments
//   - the first argument is not a struct or a reference to one
//
// FIXME: considers the first argument only.
func StructFieldsAsOptions(selector *Selector, explicitFields ...string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "struct_fields_as_options",
			run:         StructFieldsAsOptionsAction(explicitFields...),
		},
	}
}

// DisjunctionAsOptions uses the branches of the first argument's disjunction (assuming it is one) and turns them
// into options.
//
// Example:
//
//	```
//	func Panel(panel Panel|Row) {
//		this.resource.panels.append(panel)
//	}
//	```
//
// Will become:
//
//	```
//	func Panel(panel Panel) {
//		this.resource.panels.append(panel)
//	}
//
//	func Row(row Row) {
//		this.resource.panels.append(row)
//	}
//	```
//
// This action returns the option unchanged if:
//   - it has no arguments
//   - the given argument is not a disjunction or a reference to one
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

// AddAssignment adds an assignment to an existing option.
func AddAssignment(selector *Selector, assignment veneers.Assignment) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "add_assignment",
			run:         AddAssignmentAction(assignment),
		},
	}
}

// AddComments adds comments to an option.
func AddComments(selector *Selector, comments []string) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "add_comments",
			run:         AddCommentsAction(comments),
		},
	}
}

// Debug prints debugging information about an option.
func Debug(selector *Selector) Rule {
	return Rule{
		Selector: selector,
		Action: &Action{
			description: "debug",
			run:         DebugAction(),
		},
	}
}

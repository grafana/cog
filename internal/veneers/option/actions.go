package option

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type RewriteAction func(builder ast.Builder, option ast.Option) []ast.Option

func RenameAction(newName string) RewriteAction {
	return func(_ ast.Builder, option ast.Option) []ast.Option {
		newOption := option
		newOption.Name = newName

		return []ast.Option{newOption}
	}
}

// FIXME: looks at the first arg only, no way to configure that right now
func ArrayToAppendAction() RewriteAction {
	return func(_ ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 || option.Args[0].Type.Kind != ast.KindArray {
			return []ast.Option{option}
		}

		// Update the argument type from list to a single value
		oldArgs := option.Args

		newFirstArg := option.Args[0]
		newFirstArg.Type = option.Args[0].Type.AsArray().ValueType

		// Update the assignment to do an append instead of a list assignment
		oldAssignments := option.Assignments

		newFirstAssignment := option.Assignments[0]
		newFirstAssignment.Method = ast.AppendAssignment

		newOpt := option
		newOpt.Args = []ast.Argument{newFirstArg}
		newOpt.Assignments = []ast.Assignment{newFirstAssignment}

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}
		if len(oldAssignments) > 1 {
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ast.Option{newOpt}
	}
}

func OmitAction() RewriteAction {
	return func(_ ast.Builder, _ ast.Option) []ast.Option {
		return nil
	}
}

func PromoteToConstructorAction() RewriteAction {
	return func(_ ast.Builder, option ast.Option) []ast.Option {
		newOpt := option
		newOpt.IsConstructorArg = true

		return []ast.Option{newOpt}
	}
}

// FIXME: looks at the first arg only, no way to configure that right now
func StructFieldsAsArgumentsAction(explicitFields ...string) RewriteAction {
	return func(builder ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 {
			return []ast.Option{option}
		}

		firstArgType := option.Args[0].Type
		if firstArgType.Kind == ast.KindRef {
			referredObject := builder.Schema.LocateObject(firstArgType.AsRef().ReferredType)
			firstArgType = referredObject.Type
		}

		if firstArgType.Kind != ast.KindStruct {
			return []ast.Option{option}
		}

		oldArgs := option.Args
		oldAssignments := option.Assignments
		assignmentPathPrefix := oldAssignments[0].Path
		structType := firstArgType.AsStruct()

		newOpt := option
		newOpt.Args = nil
		newOpt.Assignments = nil

		for _, field := range structType.Fields {
			if explicitFields != nil && !tools.ItemInList(field.Name, explicitFields) {
				continue
			}

			var constraints []ast.TypeConstraint
			if field.Type.Kind == ast.KindScalar {
				constraints = field.Type.AsScalar().Constraints
			}

			newArg := ast.Argument{
				Name: field.Name,
				Type: field.Type,
			}

			newOpt.Args = append(newOpt.Args, newArg)

			newAssignment := ast.ArgumentAssignment(
				assignmentPathPrefix.Append(ast.PathFromStructField(field)),
				newArg,
				ast.Constraints(constraints),
			)

			newOpt.Assignments = append(newOpt.Assignments, newAssignment)
		}

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ast.Option{newOpt}
	}
}

type BooleanUnfold struct {
	OptionTrue  string
	OptionFalse string
}

func UnfoldBooleanAction(unfoldOpts BooleanUnfold) RewriteAction {
	return func(_ ast.Builder, option ast.Option) []ast.Option {
		newOpts := []ast.Option{
			{
				Name:     unfoldOpts.OptionTrue,
				Comments: option.Comments,
				Assignments: []ast.Assignment{
					ast.ConstantAssignment(option.Assignments[0].Path, true),
				},
			},

			{
				Name:     unfoldOpts.OptionFalse,
				Comments: option.Comments,
				Assignments: []ast.Assignment{
					ast.ConstantAssignment(option.Assignments[0].Path, false),
				},
			},
		}

		if option.Default != nil {
			if option.Default.ArgsValues[0].(bool) {
				newOpts[0].Default = &ast.OptionDefault{}
			} else {
				newOpts[1].Default = &ast.OptionDefault{}
			}
		}

		return newOpts
	}
}

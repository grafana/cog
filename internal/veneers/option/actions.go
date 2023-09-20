package option

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type RewriteAction func(option ast.Option) []ast.Option

func RenameAction(newName string) RewriteAction {
	return func(option ast.Option) []ast.Option {
		newOption := option
		newOption.Name = newName

		return []ast.Option{newOption}
	}
}

// FIXME: looks at the first arg only, no way to configure that right now
func ArrayToAppendAction() RewriteAction {
	return func(option ast.Option) []ast.Option {
		if len(option.Args) < 1 || option.Args[0].Type.Kind != ast.KindArray {
			return []ast.Option{option}
		}

		oldArgs := option.Args

		newFirstArg := option.Args[0]
		newFirstArg.Type = option.Args[0].Type.AsArray().ValueType

		newOpt := option
		newOpt.Args = []ast.Argument{newFirstArg}

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}

		return []ast.Option{newOpt}
	}
}

func OmitAction() RewriteAction {
	return func(_ ast.Option) []ast.Option {
		return nil
	}
}

func PromoteToConstructorAction() RewriteAction {
	return func(option ast.Option) []ast.Option {
		newOpt := option
		newOpt.IsConstructorArg = true

		return []ast.Option{newOpt}
	}
}

// FIXME: looks at the first arg only, no way to configure that right now
func StructFieldsAsArgumentsAction(explicitFields ...string) RewriteAction {
	return func(option ast.Option) []ast.Option {
		// TODO: handle the case where option.Args[0].Type is a KindRef. Follow the ref and keep working.
		if len(option.Args) < 1 || option.Args[0].Type.Kind != ast.KindStruct {
			return []ast.Option{option}
		}

		oldArgs := option.Args
		oldAssignments := option.Assignments
		assignmentPathPrefix := oldAssignments[0].Path
		structType := option.Args[0].Type.AsStruct()

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

			newOpt.Args = append(newOpt.Args, ast.Argument{
				Name: field.Name,
				Type: field.Type,
			})

			newOpt.Assignments = append(newOpt.Assignments, ast.Assignment{
				Path:              assignmentPathPrefix + "." + field.Name,
				ArgumentName:      field.Name,
				ValueType:         field.Type,
				Constraints:       constraints,
				IntoOptionalField: !field.Required,
			})
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
	return func(option ast.Option) []ast.Option {
		newOpts := []ast.Option{
			{
				Name:     unfoldOpts.OptionTrue,
				Comments: option.Comments,
				Args:     nil,
				Assignments: []ast.Assignment{
					{
						Path:              option.Assignments[0].Path,
						ValueType:         option.Assignments[0].ValueType,
						IntoOptionalField: option.Assignments[0].IntoOptionalField,
						Value:             true,
					},
				},
			},

			{
				Name:     unfoldOpts.OptionFalse,
				Comments: option.Comments,
				Args:     nil,
				Assignments: []ast.Assignment{
					{
						Path:              option.Assignments[0].Path,
						ValueType:         option.Assignments[0].ValueType,
						IntoOptionalField: option.Assignments[0].IntoOptionalField,
						Value:             false,
					},
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

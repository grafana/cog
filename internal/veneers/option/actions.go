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
		// TODO: what if there is an envelope in the value assignment?
		if newFirstAssignment.Value.Argument != nil {
			newFirstAssignment.Value.Argument.Type = newFirstArg.Type
		}

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

// FIXME: looks at the first arg only, no way to configure that right now
func DisjunctionAsOptionsAction() RewriteAction {
	return func(builder ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 {
			return []ast.Option{option}
		}

		firstArgType := option.Args[0].Type

		// "proper" disjunction
		if firstArgType.Kind == ast.KindDisjunction {
			return disjunctionAsOptions(option)
		}

		// or maybe a reference to a struct that was created to simulate a disjunction?
		if firstArgType.Kind == ast.KindRef {
			// FIXME: we only try to resolve the reference within the same package
			referredObj := builder.Schema.LocateObject(firstArgType.AsRef().ReferredType)
			// Object not found
			// TODO: LocateObject() should return a "found" boolean
			if referredObj.Name == "" {
				return []ast.Option{option}
			}

			if !referredObj.Type.IsStructGeneratedFromDisjunction() {
				return []ast.Option{option}
			}

			return disjunctionStructAsOptions(option, referredObj)
		}

		return []ast.Option{option}
	}
}

func disjunctionStructAsOptions(option ast.Option, disjunctionStruct ast.Object) []ast.Option {
	firstArgType := option.Args[0].Type
	firstAssignmentPath := option.Assignments[0].Path
	firstAssignmentMethod := option.Assignments[0].Method

	newOpts := make([]ast.Option, 0, len(disjunctionStruct.Type.AsStruct().Fields))
	for _, field := range disjunctionStruct.Type.AsStruct().Fields {
		arg := ast.Argument{Name: field.Name, Type: field.Type}

		opt := ast.Option{
			Name: field.Name,
			Args: []ast.Argument{arg},
			Assignments: []ast.Assignment{
				{
					Path: firstAssignmentPath,
					Value: ast.AssignmentValue{
						Envelope: &ast.AssignmentEnvelope{
							Type: firstArgType.AsRef(),
							Path: ast.PathFromStructField(field),
							Value: ast.AssignmentValue{
								Argument: &arg,
							},
						},
					},
					Method: firstAssignmentMethod,
				},
			},
		}

		if field.Type.Default != nil {
			opt.Default = &ast.OptionDefault{
				ArgsValues: []any{field.Type.Default},
			}
		}

		newOpts = append(newOpts, opt)
	}

	return newOpts
}

func disjunctionAsOptions(option ast.Option) []ast.Option {
	firstArgType := option.Args[0].Type
	firstAssignmentPath := option.Assignments[0].Path
	firstAssignmentMethod := option.Assignments[0].Method

	newOpts := make([]ast.Option, 0, len(firstArgType.AsDisjunction().Branches))
	for _, branch := range firstArgType.AsDisjunction().Branches {
		typeName := tools.LowerCamelCase(ast.TypeName(branch))

		arg := ast.Argument{Name: typeName, Type: branch}

		opt := ast.Option{
			Name: typeName,
			Args: []ast.Argument{arg},
			Assignments: []ast.Assignment{
				ast.ArgumentAssignment(firstAssignmentPath, arg, ast.Method(firstAssignmentMethod)),
			},
		}

		if branch.Default != nil {
			opt.Default = &ast.OptionDefault{
				ArgsValues: []any{branch.Default},
			}
		}

		newOpts = append(newOpts, opt)
	}

	return newOpts
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

package option

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers"
)

type ActionRunner func(ctx RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error)

func RenameAction(newName string) ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, option ast.Option) ([]ast.Option, error) {
		oldName := option.Name
		option.Name = newName
		option.AddToVeneerTrail(fmt.Sprintf("Rename[%s → %s]", oldName, newName))

		return []ast.Option{option}, nil
	}
}

func RenameArgumentsAction(newNames []string) ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, option ast.Option) ([]ast.Option, error) {
		if len(newNames) != len(option.Args) {
			return []ast.Option{option}, nil
		}

		for i, arg := range option.Args {
			previousName := arg.Name
			option.Args[i].Name = newNames[i]

			for j, assignment := range option.Assignments {
				if assignment.Value.Argument != nil && assignment.Value.Argument.Name == previousName {
					option.Assignments[j].Value.Argument.Name = newNames[i]
				}
			}
		}

		option.AddToVeneerTrail("RenameArguments")

		return []ast.Option{option}, nil
	}
}

func ArrayToAppendAction() ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, option ast.Option) ([]ast.Option, error) {
		if len(option.Args) != 1 || !option.Args[0].Type.IsArray() {
			return []ast.Option{option}, nil
		}

		// Update the argument type from list to a single value
		oldArgs := option.Args

		newFirstArg := option.Args[0]
		newFirstArg.Type = option.Args[0].Type.AsArray().ValueType
		newFirstArg.Name = tools.Singularize(newFirstArg.Name)

		// Update the assignment to do an append instead of a list assignment
		oldAssignments := option.Assignments

		newFirstAssignment := option.Assignments[0]
		newFirstAssignment.Method = ast.AppendAssignment
		// TODO: what if there is an envelope in the value assignment?
		if newFirstAssignment.Value.Argument != nil {
			newFirstAssignment.Value.Argument.Name = newFirstArg.Name
			newFirstAssignment.Value.Argument.Type = newFirstArg.Type
		}

		newOpt := option
		newOpt.Args = []ast.Argument{newFirstArg}
		newOpt.Assignments = []ast.Assignment{newFirstAssignment}
		newOpt.AddToVeneerTrail("ArrayToAppend")

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}
		if len(oldAssignments) > 1 {
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ast.Option{newOpt}, nil
	}
}

func MapToIndexAction() ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, option ast.Option) ([]ast.Option, error) {
		if len(option.Args) != 1 || !option.Args[0].Type.IsMap() {
			return []ast.Option{option}, nil
		}

		oldArgs := option.Args

		newFirstArg := option.Args[0]
		newFirstArg.Type = option.Args[0].Type.Map.IndexType
		newFirstArg.Name = "key"

		newSecondArg := option.Args[0]
		newSecondArg.Type = option.Args[0].Type.Map.ValueType
		newSecondArg.Name = tools.Singularize(option.Args[0].Name)

		// Update the assignment to do an append instead of a list assignment
		oldAssignments := option.Assignments

		newFirstAssignment := option.Assignments[0]
		newFirstAssignment.Method = ast.IndexAssignment
		newFirstAssignment.Path = newFirstAssignment.Path.Append(ast.Path{{
			Index: &ast.PathIndex{Argument: &newFirstArg},
			Type:  option.Args[0].Type.Map.ValueType,
		}})
		// TODO: what if there is an envelope in the value assignment?
		if newFirstAssignment.Value.Argument != nil {
			newFirstAssignment.Value.Argument.Name = newSecondArg.Name
			newFirstAssignment.Value.Argument.Type = newSecondArg.Type
		}

		newOpt := option
		newOpt.Args = []ast.Argument{newFirstArg, newSecondArg}
		newOpt.Assignments = []ast.Assignment{newFirstAssignment}
		newOpt.AddToVeneerTrail("MapToIndex")

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}
		if len(oldAssignments) > 1 {
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ast.Option{newOpt}, nil
	}
}

func OmitAction() ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, _ ast.Option) ([]ast.Option, error) {
		return nil, nil
	}
}

func VeneerTrailAsCommentsAction() ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, opt ast.Option) ([]ast.Option, error) {
		veneerTrail := tools.Map(opt.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})

		opt.Comments = append(opt.Comments, veneerTrail...)

		return []ast.Option{opt}, nil
	}
}

func StructFieldsAsArgumentsAction(explicitFields ...string) ActionRunner {
	return func(ctx RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		if len(option.Args) < 1 {
			ctx.Logger.Warn("option has no arguments: skipping transformation")
			return []ast.Option{option}, nil
		}

		firstArgType := ctx.Schemas.ResolveToType(option.Args[0].Type)

		if !firstArgType.IsStruct() {
			ctx.Logger.Warn("first argument does not resolve to a struct: skipping transformation", slog.String("type", ast.TypeName(firstArgType)))
			return []ast.Option{option}, nil
		}

		oldArgs := option.Args
		oldAssignments := option.Assignments
		assignmentPathPrefix := oldAssignments[0].Path
		structType := firstArgType.AsStruct()

		newOpt := option
		newOpt.Args = nil
		newOpt.Assignments = nil
		newOpt.Default = nil
		newOpt.AddToVeneerTrail("StructFieldsAsArguments")

		assignIntoList := assignmentPathPrefix.Last().Type.IsArray()

		newAssignments := make([]ast.Assignment, 0, len(structType.Fields))
		valuesForEnvelope := make([]ast.EnvelopeFieldValue, 0, len(structType.Fields))
		defaults := make(map[string]any)
		if option.Default != nil && len(option.Default.ArgsValues) == 1 {
			if defs, ok := option.Default.ArgsValues[0].(map[string]any); ok {
				defaults = defs
			}
		}

		for _, field := range structType.Fields {
			if explicitFields != nil && !tools.ItemInList(field.Name, explicitFields) {
				continue
			}

			var constraints []ast.TypeConstraint
			if field.Type.IsScalar() {
				constraints = field.Type.AsScalar().Constraints
			}

			// It sets the default to the args to simplify the process to extract the values in each language
			// since defaults don't have enough information to detect a reference.
			if def, ok := defaults[field.Name]; ok {
				field.Type.Default = def
			}

			newArg := ast.Argument{
				Name: field.Name,
				Type: field.Type,
			}

			// if the field has a value, it's a constant and we don't need to add it as an argument
			isConstant := field.Type.IsConcreteScalar()
			if !isConstant {
				newOpt.Args = append(newOpt.Args, newArg)
			}

			if !assignIntoList {
				var newAssignment ast.Assignment
				if isConstant {
					newAssignment = ast.ConstantAssignment(
						assignmentPathPrefix.Append(ast.PathFromStructField(field)),
						field.Type.AsScalar().Value,
					)
				} else {
					newAssignment = ast.ArgumentAssignment(
						assignmentPathPrefix.Append(ast.PathFromStructField(field)),
						newArg,
						ast.WithTypeConstraints(constraints),
						ast.Method(oldAssignments[0].Method),
					)
				}

				newAssignments = append(newAssignments, newAssignment)
			} else {
				var assignmentValue ast.AssignmentValue
				if isConstant {
					assignmentValue = ast.AssignmentValue{Constant: field.Type.AsScalar().Value}
				} else {
					assignmentValue = ast.AssignmentValue{Argument: &newArg}
				}
				valuesForEnvelope = append(valuesForEnvelope, ast.EnvelopeFieldValue{
					Path:  ast.PathFromStructField(field),
					Value: assignmentValue,
				})
			}

			if defaults[field.Name] != nil {
				if newOpt.Default == nil {
					newOpt.Default = &ast.OptionDefault{}
				}

				newOpt.Default.ArgsValues = append(newOpt.Default.ArgsValues, defaults[field.Name])
			}
		}

		if !assignIntoList {
			newOpt.Assignments = newAssignments
		} else {
			newOpt.Assignments = []ast.Assignment{
				{
					Method: ast.AppendAssignment,
					Path:   assignmentPathPrefix,
					Value: ast.AssignmentValue{
						Envelope: &ast.AssignmentEnvelope{
							Type:   assignmentPathPrefix.Last().Type.AsArray().ValueType,
							Values: valuesForEnvelope,
						},
					},
				},
			}
		}

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ast.Option{newOpt}, nil
	}
}

func StructFieldsAsOptionsAction(explicitFields ...string) ActionRunner {
	return func(ctx RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		if len(option.Args) < 1 {
			ctx.Logger.Warn("option has no arguments: skipping transformation")
			return []ast.Option{option}, nil
		}

		firstArgType := ctx.Schemas.ResolveToType(option.Args[0].Type)

		if !firstArgType.IsStruct() {
			ctx.Logger.Warn("first argument does not resolve to a struct: skipping transformation", slog.String("type", ast.TypeName(firstArgType)))
			return []ast.Option{option}, nil
		}

		var newOptions []ast.Option

		structType := firstArgType.AsStruct()
		oldAssignments := option.Assignments
		assignmentPathPrefix := oldAssignments[0].Path

		for _, field := range structType.Fields {
			if explicitFields != nil && !tools.ItemInList(field.Name, explicitFields) {
				continue
			}

			newOpt := ast.Option{
				Name:     field.Name,
				Comments: field.Comments,
				Args: []ast.Argument{
					{Name: field.Name, Type: field.Type},
				},
				Assignments: []ast.Assignment{
					ast.FieldAssignment(field),
				},
			}
			newOpt.AddToVeneerTrail("StructFieldsAsOptions")

			newOpt.Assignments[0].Path = assignmentPathPrefix.Append(newOpt.Assignments[0].Path)

			if field.Type.Default != nil {
				newOpt.Default = &ast.OptionDefault{
					ArgsValues: []any{field.Type.Default},
				}
			}

			newOptions = append(newOptions, newOpt)
		}

		return newOptions, nil
	}
}
func DisjunctionAsOptionsAction(argumentIndex int) ActionRunner {
	return func(ctx RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		if len(option.Args) == 0 {
			ctx.Logger.Warn("option has no arguments: skipping transformation")
			return []ast.Option{option}, nil
		}

		targetArgType := option.Args[argumentIndex].Type

		// "proper" disjunction
		if targetArgType.IsDisjunction() {
			return disjunctionAsOptions(option, argumentIndex)
		}

		// or maybe a reference to a struct that was created to simulate a disjunction?
		if targetArgType.IsRef() {
			referredType := ctx.Schemas.ResolveToType(targetArgType)
			if !referredType.IsStructGeneratedFromDisjunction() {
				return []ast.Option{option}, nil
			}

			return disjunctionStructAsOptions(option, referredType, argumentIndex)
		}

		return []ast.Option{option}, nil
	}
}

func disjunctionStructAsOptions(option ast.Option, disjunctionStruct ast.Type, argIndex int) ([]ast.Option, error) {
	newOpts := make([]ast.Option, 0, len(disjunctionStruct.AsStruct().Fields))
	for _, field := range disjunctionStruct.AsStruct().Fields {
		optClone := option.DeepCopy()

		arg := ast.Argument{Name: field.Name, Type: field.Type}
		args := optClone.Args[0:argIndex]
		args = append(args, arg)
		if len(option.Args) > argIndex+1 {
			args = append(args, option.Args[argIndex+1:]...)
		}

		assignments := optClone.Assignments
		for i, assignment := range assignments {
			if assignment.Value.Argument == nil || assignment.Value.Argument.Name != option.Args[argIndex].Name {
				continue
			}

			assignments[i] = ast.Assignment{
				Path: assignments[i].Path,
				Value: ast.AssignmentValue{
					Envelope: &ast.AssignmentEnvelope{
						Type: option.Args[argIndex].Type,
						Values: []ast.EnvelopeFieldValue{
							{
								Path:  ast.PathFromStructField(field),
								Value: ast.AssignmentValue{Argument: &arg},
							},
						},
					},
				},
				Method: assignments[i].Method,
			}
			break
		}

		opt := ast.Option{
			Name:        field.Name,
			Args:        args,
			Assignments: assignments,
		}
		opt.AddToVeneerTrail("DisjunctionAsOptions")

		if field.Type.Default != nil {
			opt.Default = &ast.OptionDefault{
				ArgsValues: []any{field.Type.Default},
			}
		}

		newOpts = append(newOpts, opt)
	}

	return newOpts, nil
}

func disjunctionAsOptions(option ast.Option, argIndex int) ([]ast.Option, error) {
	disjunction := option.Args[argIndex].Type.AsDisjunction()

	newOpts := make([]ast.Option, 0, len(disjunction.Branches))
	for _, branch := range disjunction.Branches {
		optClone := option.DeepCopy()
		typeName := tools.LowerCamelCase(ast.TypeName(branch))

		arg := ast.Argument{Name: typeName, Type: branch}

		args := optClone.Args[0:argIndex]
		args = append(args, arg)
		if len(option.Args) > argIndex+1 {
			args = append(args, option.Args[argIndex+1:]...)
		}

		assignments := optClone.Assignments
		for i, assignment := range assignments {
			if assignment.Value.Argument == nil || assignment.Value.Argument.Name != option.Args[argIndex].Name {
				continue
			}

			assignments[i] = ast.ArgumentAssignment(
				assignments[i].Path,
				arg,
				ast.Method(assignments[i].Method),
			)
			break
		}

		opt := ast.Option{
			Name:        typeName,
			Args:        args,
			Assignments: assignments,
		}
		opt.AddToVeneerTrail("DisjunctionAsOptions")

		if branch.Default != nil {
			opt.Default = &ast.OptionDefault{
				ArgsValues: []any{branch.Default},
			}
		}

		newOpts = append(newOpts, opt)
	}

	return newOpts, nil
}

type BooleanUnfold struct {
	OptionTrue  string
	OptionFalse string
}

func UnfoldBooleanAction(unfoldOpts BooleanUnfold) ActionRunner {
	return func(_ RuleCtx, _ ast.Builder, option ast.Option) ([]ast.Option, error) {
		intoType := option.Assignments[0].Path.Last().Type

		if !intoType.IsScalar() || intoType.Scalar.ScalarKind != ast.KindBool {
			return []ast.Option{option}, nil
		}

		newOpts := []ast.Option{
			{
				Name:     unfoldOpts.OptionTrue,
				Comments: option.Comments,
				Assignments: []ast.Assignment{
					ast.ConstantAssignment(option.Assignments[0].Path, true),
				},
				VeneerTrail: append([]string{}, option.VeneerTrail...),
			},

			{
				Name:     unfoldOpts.OptionFalse,
				Comments: option.Comments,
				Assignments: []ast.Assignment{
					ast.ConstantAssignment(option.Assignments[0].Path, false),
				},
				VeneerTrail: append([]string{}, option.VeneerTrail...),
			},
		}

		if option.Default != nil {
			if val, ok := option.Default.ArgsValues[0].(bool); ok && val {
				newOpts[0].Default = &ast.OptionDefault{}
			} else {
				newOpts[1].Default = &ast.OptionDefault{}
			}
		}

		newOpts[0].AddToVeneerTrail("UnfoldBoolean")
		newOpts[1].AddToVeneerTrail("UnfoldBoolean")

		return newOpts, nil
	}
}

func DuplicateAction(duplicateName string) ActionRunner {
	return func(_ RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		duplicateOpt := option.DeepCopy()
		duplicateOpt.Name = duplicateName
		duplicateOpt.AddToVeneerTrail(fmt.Sprintf("Duplicate[%s]", option.Name))

		return []ast.Option{option, duplicateOpt}, nil
	}
}

func AddAssignmentAction(assignment veneers.Assignment) ActionRunner {
	return func(ctx RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		irAssignment, err := assignment.AsIR(ctx.Schemas, builder)
		if err != nil {
			return nil, err
		}

		option.Assignments = append(option.Assignments, irAssignment)
		option.AddToVeneerTrail(fmt.Sprintf("AddAssignment[%s]", irAssignment.Path.String()))

		return []ast.Option{option}, nil
	}
}

func AddCommentsAction(comments []string) ActionRunner {
	return func(_ RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		option.Comments = append(option.Comments, comments...)
		option.AddToVeneerTrail(fmt.Sprintf("AddComments[%s]", strings.Join(comments, " ")))

		return []ast.Option{option}, nil
	}
}

func DebugAction() ActionRunner {
	return func(_ RuleCtx, builder ast.Builder, option ast.Option) ([]ast.Option, error) {
		marshaled, err := json.MarshalIndent(option, "", "  ")
		if err != nil {
			return nil, err
		}

		fmt.Printf("[debug] option %s.%s.%s:\n", builder.Package, builder.Name, option.Name)
		fmt.Println(string(marshaled))

		return []ast.Option{option}, nil
	}
}

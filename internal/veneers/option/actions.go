package option

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers"
	"github.com/grafana/cog/pkg/ir"
)

type ActionRunner func(ctx RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error)

func RenameAction(newName string) ActionRunner {
	return func(_ RuleCtx, _ ir.Builder, option ir.Option) ([]ir.Option, error) {
		oldName := option.Name
		option.Name = newName
		option.AddToVeneerTrail(fmt.Sprintf("Rename[%s → %s]", oldName, newName))

		return []ir.Option{option}, nil
	}
}

func RenameArgumentsAction(newNames []string) ActionRunner {
	return func(ctx RuleCtx, _ ir.Builder, option ir.Option) ([]ir.Option, error) {
		if len(newNames) != len(option.Args) {
			ctx.Logger.Warn("the number of new argument names does not match the number of arguments: skipping transformation", slog.Int("new_names_count", len(newNames)), slog.Int("args_count", len(option.Args)))
			return []ir.Option{option}, nil
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

		return []ir.Option{option}, nil
	}
}

func ArrayToAppendAction() ActionRunner {
	return func(ctx RuleCtx, _ ir.Builder, option ir.Option) ([]ir.Option, error) {
		if len(option.Args) != 1 {
			ctx.Logger.Warn("expecting a single argument: skipping transformation", slog.Int("args_count", len(option.Args)))
			return []ir.Option{option}, nil
		}

		if !option.Args[0].Type.IsArray() {
			ctx.Logger.Warn("first argument is not an array: skipping transformation", slog.String("type", ir.TypeName(option.Args[0].Type)))
			return []ir.Option{option}, nil
		}

		// Update the argument type from list to a single value
		oldArgs := option.Args

		newFirstArg := option.Args[0]
		newFirstArg.Type = option.Args[0].Type.AsArray().ValueType
		newFirstArg.Name = tools.Singularize(newFirstArg.Name)

		// Update the assignment to do an append instead of a list assignment
		oldAssignments := option.Assignments

		newFirstAssignment := option.Assignments[0]
		newFirstAssignment.Method = ir.AppendAssignment
		// TODO: what if there is an envelope in the value assignment?
		if newFirstAssignment.Value.Argument != nil {
			newFirstAssignment.Value.Argument.Name = newFirstArg.Name
			newFirstAssignment.Value.Argument.Type = newFirstArg.Type
		}

		newOpt := option
		newOpt.Args = []ir.Argument{newFirstArg}
		newOpt.Assignments = []ir.Assignment{newFirstAssignment}
		newOpt.AddToVeneerTrail("ArrayToAppend")

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}
		if len(oldAssignments) > 1 {
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ir.Option{newOpt}, nil
	}
}

func MapToIndexAction() ActionRunner {
	return func(ctx RuleCtx, _ ir.Builder, option ir.Option) ([]ir.Option, error) {
		if len(option.Args) != 1 {
			ctx.Logger.Warn("expecting a single argument: skipping transformation", slog.Int("args_count", len(option.Args)))
			return []ir.Option{option}, nil
		}

		if !option.Args[0].Type.IsMap() {
			ctx.Logger.Warn("first argument is not a map: skipping transformation", slog.String("type", ir.TypeName(option.Args[0].Type)))
			return []ir.Option{option}, nil
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
		newFirstAssignment.Method = ir.IndexAssignment
		newFirstAssignment.Path = newFirstAssignment.Path.Append(ir.Path{{
			Index: &ir.PathIndex{Argument: &newFirstArg},
			Type:  option.Args[0].Type.Map.ValueType,
		}})
		// TODO: what if there is an envelope in the value assignment?
		if newFirstAssignment.Value.Argument != nil {
			newFirstAssignment.Value.Argument.Name = newSecondArg.Name
			newFirstAssignment.Value.Argument.Type = newSecondArg.Type
		}

		newOpt := option
		newOpt.Args = []ir.Argument{newFirstArg, newSecondArg}
		newOpt.Assignments = []ir.Assignment{newFirstAssignment}
		newOpt.AddToVeneerTrail("MapToIndex")

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}
		if len(oldAssignments) > 1 {
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ir.Option{newOpt}, nil
	}
}

func OmitAction() ActionRunner {
	return func(_ RuleCtx, _ ir.Builder, _ ir.Option) ([]ir.Option, error) {
		return nil, nil
	}
}

func VeneerTrailAsCommentsAction() ActionRunner {
	return func(_ RuleCtx, _ ir.Builder, opt ir.Option) ([]ir.Option, error) {
		veneerTrail := tools.Map(opt.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})

		opt.Comments = append(opt.Comments, veneerTrail...)

		return []ir.Option{opt}, nil
	}
}

func StructFieldsAsArgumentsAction(explicitFields ...string) ActionRunner {
	return func(ctx RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		if len(option.Args) == 0 {
			ctx.Logger.Warn("option has no arguments: skipping transformation")
			return []ir.Option{option}, nil
		}

		firstArgType := ctx.Schemas.Resolve(option.Args[0].Type)
		if !firstArgType.IsStruct() {
			ctx.Logger.Warn("first argument does not resolve to a struct: skipping transformation", slog.String("type", ir.TypeName(firstArgType)))
			return []ir.Option{option}, nil
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

		newAssignments := make([]ir.Assignment, 0, len(structType.Fields))
		valuesForEnvelope := make([]ir.EnvelopeFieldValue, 0, len(structType.Fields))
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

			var constraints []ir.TypeConstraint
			if field.Type.IsScalar() {
				constraints = field.Type.AsScalar().Constraints
			}

			// It sets the default to the args to simplify the process to extract the values in each language
			// since defaults don't have enough information to detect a reference.
			if def, ok := defaults[field.Name]; ok {
				field.Type.Default = def
			}

			newArg := ir.Argument{
				Name: field.Name,
				Type: field.Type,
			}

			// if the field has a value, it's a constant and we don't need to add it as an argument
			isConstant := field.Type.IsConcreteScalar()
			if !isConstant {
				newOpt.Args = append(newOpt.Args, newArg)
			}

			if !assignIntoList {
				var newAssignment ir.Assignment
				if isConstant {
					newAssignment = ir.ConstantAssignment(
						assignmentPathPrefix.Append(ir.PathFromStructField(field)),
						field.Type.AsScalar().Value,
					)
				} else {
					newAssignment = ir.ArgumentAssignment(
						assignmentPathPrefix.Append(ir.PathFromStructField(field)),
						newArg,
						ir.WithTypeConstraints(constraints),
						ir.Method(oldAssignments[0].Method),
					)
				}

				newAssignments = append(newAssignments, newAssignment)
			} else {
				var assignmentValue ir.AssignmentValue
				if isConstant {
					assignmentValue = ir.AssignmentValue{Constant: field.Type.AsScalar().Value}
				} else {
					assignmentValue = ir.AssignmentValue{Argument: &newArg}
				}
				valuesForEnvelope = append(valuesForEnvelope, ir.EnvelopeFieldValue{
					Path:  ir.PathFromStructField(field),
					Value: assignmentValue,
				})
			}

			if defaults[field.Name] != nil {
				if newOpt.Default == nil {
					newOpt.Default = &ir.OptionDefault{}
				}

				newOpt.Default.ArgsValues = append(newOpt.Default.ArgsValues, defaults[field.Name])
			}
		}

		if !assignIntoList {
			newOpt.Assignments = newAssignments
		} else {
			newOpt.Assignments = []ir.Assignment{
				{
					Method: ir.AppendAssignment,
					Path:   assignmentPathPrefix,
					Value: ir.AssignmentValue{
						Envelope: &ir.AssignmentEnvelope{
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

		return []ir.Option{newOpt}, nil
	}
}

func StructFieldsAsOptionsAction(explicitFields ...string) ActionRunner {
	return func(ctx RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		if len(option.Args) == 0 {
			ctx.Logger.Warn("option has no arguments: skipping transformation")
			return []ir.Option{option}, nil
		}

		firstArgType := ctx.Schemas.Resolve(option.Args[0].Type)
		if !firstArgType.IsStruct() {
			ctx.Logger.Warn("first argument does not resolve to a struct: skipping transformation", slog.String("type", ir.TypeName(firstArgType)))
			return []ir.Option{option}, nil
		}

		var newOptions []ir.Option

		structType := firstArgType.AsStruct()
		oldAssignments := option.Assignments
		assignmentPathPrefix := oldAssignments[0].Path

		for _, field := range structType.Fields {
			if explicitFields != nil && !tools.ItemInList(field.Name, explicitFields) {
				continue
			}

			newOpt := ir.Option{
				Name:     field.Name,
				Comments: field.Comments,
				Args: []ir.Argument{
					{Name: field.Name, Type: field.Type},
				},
				Assignments: []ir.Assignment{
					ir.FieldAssignment(field),
				},
			}
			newOpt.AddToVeneerTrail("StructFieldsAsOptions")

			newOpt.Assignments[0].Path = assignmentPathPrefix.Append(newOpt.Assignments[0].Path)

			if field.Type.Default != nil {
				newOpt.Default = &ir.OptionDefault{
					ArgsValues: []any{field.Type.Default},
				}
			}

			newOptions = append(newOptions, newOpt)
		}

		return newOptions, nil
	}
}
func DisjunctionAsOptionsAction(argumentIndex int) ActionRunner {
	return func(ctx RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		if len(option.Args) == 0 {
			ctx.Logger.Warn("option has no arguments: skipping transformation")
			return []ir.Option{option}, nil
		}

		targetArgType := option.Args[argumentIndex].Type

		// "proper" disjunction
		if targetArgType.IsDisjunction() {
			return disjunctionAsOptions(option, argumentIndex)
		}

		// or maybe a reference to a struct that was created to simulate a disjunction?
		if targetArgType.IsRef() {
			referredType := ctx.Schemas.Resolve(targetArgType)
			if !referredType.IsStructGeneratedFromDisjunction() {
				ctx.Logger.Warn("argument is not a ref to a disjunction: skipping transformation", slog.String("type", ir.TypeName(referredType)))
				return []ir.Option{option}, nil
			}

			return disjunctionStructAsOptions(option, referredType, argumentIndex)
		}

		ctx.Logger.Warn("argument is not a disjunction: skipping transformation", slog.String("type", ir.TypeName(targetArgType)))
		return []ir.Option{option}, nil
	}
}

func disjunctionStructAsOptions(option ir.Option, disjunctionStruct ir.Type, argIndex int) ([]ir.Option, error) {
	newOpts := make([]ir.Option, 0, len(disjunctionStruct.AsStruct().Fields))
	for _, field := range disjunctionStruct.AsStruct().Fields {
		optClone := option.DeepCopy()

		arg := ir.Argument{Name: field.Name, Type: field.Type}
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

			assignments[i] = ir.Assignment{
				Path: assignments[i].Path,
				Value: ir.AssignmentValue{
					Envelope: &ir.AssignmentEnvelope{
						Type: option.Args[argIndex].Type,
						Values: []ir.EnvelopeFieldValue{
							{
								Path:  ir.PathFromStructField(field),
								Value: ir.AssignmentValue{Argument: &arg},
							},
						},
					},
				},
				Method: assignments[i].Method,
			}
			break
		}

		opt := ir.Option{
			Name:        field.Name,
			Args:        args,
			Assignments: assignments,
		}
		opt.AddToVeneerTrail("DisjunctionAsOptions")

		if field.Type.Default != nil {
			opt.Default = &ir.OptionDefault{
				ArgsValues: []any{field.Type.Default},
			}
		}

		newOpts = append(newOpts, opt)
	}

	return newOpts, nil
}

func disjunctionAsOptions(option ir.Option, argIndex int) ([]ir.Option, error) {
	disjunction := option.Args[argIndex].Type.AsDisjunction()

	newOpts := make([]ir.Option, 0, len(disjunction.Branches))
	for _, branch := range disjunction.Branches {
		optClone := option.DeepCopy()
		typeName := tools.LowerCamelCase(ir.TypeName(branch))

		arg := ir.Argument{Name: typeName, Type: branch}

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

			assignments[i] = ir.ArgumentAssignment(
				assignments[i].Path,
				arg,
				ir.Method(assignments[i].Method),
			)
			break
		}

		opt := ir.Option{
			Name:        typeName,
			Args:        args,
			Assignments: assignments,
		}
		opt.AddToVeneerTrail("DisjunctionAsOptions")

		if branch.Default != nil {
			opt.Default = &ir.OptionDefault{
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
	return func(ctx RuleCtx, _ ir.Builder, option ir.Option) ([]ir.Option, error) {
		intoType := option.Assignments[0].Path.Last().Type

		if !intoType.IsScalar() || intoType.Scalar.ScalarKind != ir.KindBool {
			ctx.Logger.Warn("first assignment is not an boolean: skipping transformation", slog.String("type", ir.TypeName(intoType)))
			return []ir.Option{option}, nil
		}

		newOpts := []ir.Option{
			{
				Name:     unfoldOpts.OptionTrue,
				Comments: option.Comments,
				Assignments: []ir.Assignment{
					ir.ConstantAssignment(option.Assignments[0].Path, true),
				},
				VeneerTrail: append([]string{}, option.VeneerTrail...),
			},

			{
				Name:     unfoldOpts.OptionFalse,
				Comments: option.Comments,
				Assignments: []ir.Assignment{
					ir.ConstantAssignment(option.Assignments[0].Path, false),
				},
				VeneerTrail: append([]string{}, option.VeneerTrail...),
			},
		}

		if option.Default != nil {
			if val, ok := option.Default.ArgsValues[0].(bool); ok && val {
				newOpts[0].Default = &ir.OptionDefault{}
			} else {
				newOpts[1].Default = &ir.OptionDefault{}
			}
		}

		newOpts[0].AddToVeneerTrail("UnfoldBoolean")
		newOpts[1].AddToVeneerTrail("UnfoldBoolean")

		return newOpts, nil
	}
}

func DuplicateAction(duplicateName string) ActionRunner {
	return func(_ RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		duplicateOpt := option.DeepCopy()
		duplicateOpt.Name = duplicateName
		duplicateOpt.AddToVeneerTrail(fmt.Sprintf("Duplicate[%s]", option.Name))

		return []ir.Option{option, duplicateOpt}, nil
	}
}

func AddAssignmentAction(assignment veneers.Assignment) ActionRunner {
	return func(ctx RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		irAssignment, err := assignment.AsIR(ctx.Schemas, builder)
		if err != nil {
			return nil, err
		}

		option.Assignments = append(option.Assignments, irAssignment)
		option.AddToVeneerTrail(fmt.Sprintf("AddAssignment[%s]", irAssignment.Path.String()))

		return []ir.Option{option}, nil
	}
}

func AddCommentsAction(comments []string) ActionRunner {
	return func(_ RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		option.Comments = append(option.Comments, comments...)
		option.AddToVeneerTrail(fmt.Sprintf("AddComments[%s]", strings.Join(comments, " ")))

		return []ir.Option{option}, nil
	}
}

func DebugAction() ActionRunner {
	return func(_ RuleCtx, builder ir.Builder, option ir.Option) ([]ir.Option, error) {
		marshaled, err := json.MarshalIndent(option, "", "  ")
		if err != nil {
			return nil, err
		}

		fmt.Printf("[debug] option %s.%s.%s:\n", builder.Package, builder.Name, option.Name)
		fmt.Println(string(marshaled))

		return []ir.Option{option}, nil
	}
}

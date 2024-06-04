package option

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers"
)

type RewriteAction func(schemas ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option

// RenameAction renames an option.
func RenameAction(newName string) RewriteAction {
	return func(_ ast.Schemas, _ ast.Builder, option ast.Option) []ast.Option {
		oldName := option.Name
		option.Name = newName
		option.AddToVeneerTrail(fmt.Sprintf("Rename[%s → %s]", oldName, newName))

		return []ast.Option{option}
	}
}

// ArrayToAppendAction updates the option to perform an "append" assignment.
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
//   - it has no arguments
//   - the argument is not an array
//
// FIXME: considers the first arg only.
func ArrayToAppendAction() RewriteAction {
	return func(_ ast.Schemas, _ ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 || !option.Args[0].Type.IsArray() {
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
		newOpt.AddToVeneerTrail("ArrayToAppend")

		if len(oldArgs) > 1 {
			newOpt.Args = append(newOpt.Args, oldArgs[1:]...)
		}
		if len(oldAssignments) > 1 {
			newOpt.Assignments = append(newOpt.Assignments, oldAssignments[1:]...)
		}

		return []ast.Option{newOpt}
	}
}

// OmitAction removes an option.
func OmitAction() RewriteAction {
	return func(_ ast.Schemas, _ ast.Builder, _ ast.Option) []ast.Option {
		return nil
	}
}

// VeneerTrailAsCommentsAction removes an option.
func VeneerTrailAsCommentsAction() RewriteAction {
	return func(_ ast.Schemas, _ ast.Builder, opt ast.Option) []ast.Option {
		veneerTrail := tools.Map(opt.VeneerTrail, func(veneer string) string {
			return fmt.Sprintf("Modified by veneer '%s'", veneer)
		})

		opt.Comments = append(opt.Comments, veneerTrail...)

		return []ast.Option{opt}
	}
}

// StructFieldsAsArgumentsAction uses the fields of the first argument's struct (assuming it is one) and turns them
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
func StructFieldsAsArgumentsAction(explicitFields ...string) RewriteAction {
	return func(schemas ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 {
			return []ast.Option{option}
		}

		firstArgType := option.Args[0].Type
		if firstArgType.IsRef() {
			referredObject, found := schemas.LocateObject(firstArgType.AsRef().ReferredPkg, firstArgType.AsRef().ReferredType)
			if found {
				firstArgType = referredObject.Type
			}
		}

		if !firstArgType.IsStruct() {
			return []ast.Option{option}
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

		return []ast.Option{newOpt}
	}
}

// StructFieldsAsOptionsAction uses the fields of the first argument's struct (assuming it is one) and turns them
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
func StructFieldsAsOptionsAction(explicitFields ...string) RewriteAction {
	return func(schemas ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 {
			return []ast.Option{option}
		}

		firstArgType := option.Args[0].Type
		if firstArgType.IsRef() {
			referredObject, found := schemas.LocateObject(firstArgType.AsRef().ReferredPkg, firstArgType.AsRef().ReferredType)
			if found {
				firstArgType = referredObject.Type
			}
		}

		if !firstArgType.IsStruct() {
			return []ast.Option{option}
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

		return newOptions
	}
}

// DisjunctionAsOptionsAction uses the branches of the first argument's disjunction (assuming it is one) and turns them
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
//   - the first argument is not a disjunction or a reference to one
//
// FIXME: considers the first argument only.
func DisjunctionAsOptionsAction() RewriteAction {
	return func(schemas ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option {
		if len(option.Args) < 1 {
			return []ast.Option{option}
		}

		firstArgType := option.Args[0].Type

		// "proper" disjunction
		if firstArgType.IsDisjunction() {
			return disjunctionAsOptions(option)
		}

		// or maybe a reference to a struct that was created to simulate a disjunction?
		if firstArgType.IsRef() {
			// FIXME: we only try to resolve the reference within the same package
			referredObj, found := schemas.LocateObject(firstArgType.AsRef().ReferredPkg, firstArgType.AsRef().ReferredType)
			if !found {
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
							Type: firstArgType,
							Values: []ast.EnvelopeFieldValue{
								{
									Path:  ast.PathFromStructField(field),
									Value: ast.AssignmentValue{Argument: &arg},
								},
							},
						},
					},
					Method: firstAssignmentMethod,
				},
			},
		}
		opt.AddToVeneerTrail("DisjunctionAsOptions")

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
		opt.AddToVeneerTrail("DisjunctionAsOptions")

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

// UnfoldBooleanAction transforms an option accepting a boolean argument into two argument-less options.
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
func UnfoldBooleanAction(unfoldOpts BooleanUnfold) RewriteAction {
	return func(_ ast.Schemas, _ ast.Builder, option ast.Option) []ast.Option {
		intoType := option.Assignments[0].Path.Last().Type

		if !intoType.IsScalar() || intoType.Scalar.ScalarKind != ast.KindBool {
			return []ast.Option{option}
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

		return newOpts
	}
}

func DuplicateAction(duplicateName string) RewriteAction {
	return func(_ ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option {
		duplicateOpt := option.DeepCopy()
		duplicateOpt.Name = duplicateName
		duplicateOpt.AddToVeneerTrail(fmt.Sprintf("Duplicate[%s]", option.Name))

		return []ast.Option{option, duplicateOpt}
	}
}

// AddAssignmentAction adds an assignment to an existing option.
func AddAssignmentAction(assignment veneers.Assignment) RewriteAction {
	return func(_ ast.Schemas, builder ast.Builder, option ast.Option) []ast.Option {
		irAssignment, err := assignment.AsIR(ast.Builders{builder}, builder)
		if err != nil {
			// TODO: let veneers return errors
			option.AddToVeneerTrail(fmt.Sprintf("AddAssignment[err=%s]", err.Error()))
			return []ast.Option{option}
		}

		option.Assignments = append(option.Assignments, irAssignment)
		option.AddToVeneerTrail(fmt.Sprintf("AddAssignment[%s]", irAssignment.Path.String()))

		return []ast.Option{option}
	}
}

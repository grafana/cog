package builder

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers"
)

type RewriteRule func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error)

func mapToSelected(selector Selector, mapFunc func(builders ast.Builders, builder ast.Builder) (ast.Builder, error)) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, b := range builders {
			if !selector(schemas, b) {
				continue
			}

			newBuilder, err := mapFunc(builders, b)
			if err != nil {
				return nil, err
			}

			builders[i] = newBuilder
		}

		return builders, nil
	}
}

func mergeOptions(fromBuilder ast.Builder, intoBuilder ast.Builder, underPath ast.Path, excludeOptions []string, renameOptions map[string]string) (ast.Builder, error) {
	newBuilder := intoBuilder

	if renameOptions == nil {
		renameOptions = map[string]string{}
	}

	for _, opt := range fromBuilder.Options {
		if tools.ItemInList(opt.Name, excludeOptions) {
			continue
		}

		newOpt := opt
		newOpt.Assignments = nil

		if as, found := renameOptions[newOpt.Name]; found {
			newOpt.Name = as
		}

		for _, assignment := range opt.Assignments {
			newAssignment := assignment
			newAssignment.Path = underPath.Append(assignment.Path)

			newOpt.Assignments = append(newOpt.Assignments, newAssignment)
		}

		newBuilder.Options = append(newBuilder.Options, newOpt)
	}

	return newBuilder, nil
}

func Omit(selector Selector) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		filteredBuilders := make([]ast.Builder, 0, len(builders))

		for _, builder := range builders {
			if selector(schemas, builder) {
				continue
			}

			filteredBuilders = append(filteredBuilders, builder)
		}

		return filteredBuilders, nil
	}
}

func MergeInto(selector Selector, sourceBuilderName string, underPath string, excludeOptions []string, renameOptions map[string]string) RewriteRule {
	return mapToSelected(selector, func(builders ast.Builders, destinationBuilder ast.Builder) (ast.Builder, error) {
		sourceBuilder, found := builders.LocateByObject(destinationBuilder.For.SelfRef.ReferredPkg, sourceBuilderName)
		if !found {
			// We couldn't find the source builder: let's return the selected builder untouched.
			return destinationBuilder, nil
		}

		newRoot, err := destinationBuilder.MakePath(builders, underPath)
		if err != nil {
			return destinationBuilder, err
		}

		// TODO: initializations
		newBuilder, err := mergeOptions(sourceBuilder, destinationBuilder, newRoot, excludeOptions, renameOptions)
		if err != nil {
			return ast.Builder{}, err
		}

		newBuilder.AddToVeneerTrail("MergeInto")

		return newBuilder, nil
	})
}

func composePanelType(builders ast.Builders, panelType string, panelBuilder ast.Builder, composableBuilders ast.Builders, panelOptionsToExclude []string) (ast.Builder, error) {
	newBuilder := ast.Builder{
		Package:     composableBuilders[0].Package,
		For:         panelBuilder.For,
		Name:        panelBuilder.For.Name,
		Constructor: panelBuilder.Constructor,
		Properties:  panelBuilder.Properties,
	}

	typeField, ok := panelBuilder.For.Type.AsStruct().FieldByName("type")
	if !ok {
		return panelBuilder, fmt.Errorf("could not find field 'type' in panel builder")
	}

	typeAssignment := ast.ConstantAssignment(ast.PathFromStructField(typeField), panelType)
	newBuilder.Constructor.Assignments = append(newBuilder.Constructor.Assignments, typeAssignment)

	// re-add panel-related options
	for _, panelOpt := range panelBuilder.Options {
		// this value is a constant that depends on the plugin being composed into a panel
		if panelOpt.Name == "type" {
			continue
		}

		// We don't need these options anymore since we're composing them.
		if panelOpt.Name == "options" || panelOpt.Name == "custom" {
			continue
		}

		// Is the option explicitly excluded?
		if tools.ItemInList(panelOpt.Name, panelOptionsToExclude) {
			continue
		}

		newBuilder.Options = append(newBuilder.Options, panelOpt)
	}

	compositionMap := map[string]string{ // Builder → assignment root path
		"Options":     "options",
		"FieldConfig": "fieldConfig.defaults.custom",
	}

	for _, composableBuilder := range composableBuilders {
		underPath, exists := compositionMap[composableBuilder.For.Name]
		if !exists {
			// schemas for composable panels can define more types than just "Options"
			// and "FieldConfig": we need to leave these objects untouched and
			// compose only the builders that we know of.
			continue
		}

		newRoot, err := newBuilder.MakePath(builders, underPath)
		if err != nil {
			return newBuilder, err
		}

		ref := composableBuilder.For.SelfRef
		refType := ast.NewRef(ref.ReferredPkg, ref.ReferredType)
		newRoot[len(newRoot)-1].TypeHint = &refType

		newBuilder, err = mergeOptions(composableBuilder, newBuilder, newRoot, nil, nil)
		if err != nil {
			return newBuilder, err
		}
	}

	return newBuilder, nil
}

func ComposeDashboardPanel(selector Selector, panelBuilderName string, panelOptionsToExclude []string) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		panelBuilderPkg, panelBuilderNameWithoutPkg, found := strings.Cut(panelBuilderName, ".")
		if !found {
			return nil, fmt.Errorf("panelBuilderName '%s' is incorrect: no package found", panelBuilderPkg)
		}

		panelBuilder, found := builders.LocateByObject(panelBuilderPkg, panelBuilderNameWithoutPkg)
		if !found {
			// We couldn't find the panel builder: let's return the builders untouched.
			return builders, nil
		}

		// - add to newBuilders all the builders that are not composable (ie: don't comply to the selector)
		// - build a map of composable builders, indexed by panel type
		// - aggregate the composable builders into a new, composed panel builder
		// - add the new composed panel builders to newBuilders

		newBuilders := make([]ast.Builder, 0, len(builders))
		composableBuilders := make(map[string]ast.Builders)

		for _, builder := range builders {
			// the builder is for a composable type
			if selector(schemas, builder) {
				schema, found := schemas.Locate(builder.For.SelfRef.ReferredPkg)
				if !found {
					continue
				}

				panelType := schema.Metadata.Identifier
				composableBuilders[panelType] = append(composableBuilders[panelType], builder)
				continue
			}

			newBuilders = append(newBuilders, builder)
		}

		for panelType, buildersForType := range composableBuilders {
			composedBuilder, err := composePanelType(builders, panelType, panelBuilder, buildersForType, panelOptionsToExclude)
			if err != nil {
				return nil, err
			}

			composedBuilder.AddToVeneerTrail("ComposeDashboardPanel")

			newBuilders = append(newBuilders, composedBuilder)
		}

		return newBuilders, nil
	}
}

func Rename(selector Selector, newName string) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			builders[i].Name = newName
			builders[i].AddToVeneerTrail("Rename")
		}

		return builders, nil
	}
}

func VeneerTrailAsComments(selector Selector) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			veneerTrail := tools.Map(builder.VeneerTrail, func(veneer string) string {
				return fmt.Sprintf("Modified by veneer '%s'", veneer)
			})

			builders[i].For.Comments = append(builders[i].For.Comments, veneerTrail...)
		}

		return builders, nil
	}
}

func Properties(selector Selector, properties []ast.StructField) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			builders[i].Properties = append(builders[i].Properties, properties...)
			builders[i].AddToVeneerTrail("Properties")
		}

		return builders, nil
	}
}

func Duplicate(selector Selector, duplicateName string, excludeOptions []string) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		var newBuilders ast.Builders

		for _, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			duplicatedBuilder := builder.DeepCopy()
			duplicatedBuilder.Name = duplicateName
			duplicatedBuilder.AddToVeneerTrail(fmt.Sprintf("Duplicate[%s.%s]", builder.Package, builder.Name))

			if len(excludeOptions) != 0 {
				duplicatedBuilder.Options = tools.Filter(duplicatedBuilder.Options, func(option ast.Option) bool {
					return !tools.StringInListEqualFold(option.Name, excludeOptions)
				})
			}

			newBuilders = append(newBuilders, duplicatedBuilder)
		}

		return append(builders, newBuilders...), nil
	}
}

type Initialization struct {
	PropertyPath string
	Value        any
}

func Initialize(selector Selector, statements []Initialization) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			veneerDebug := make([]string, 0, len(statements))
			for _, statement := range statements {
				path, err := builders[i].MakePath(builders, statement.PropertyPath)
				if err != nil {
					return nil, err
				}

				builders[i].Constructor.Assignments = append(builders[i].Constructor.Assignments, ast.ConstantAssignment(path, statement.Value))
				veneerDebug = append(veneerDebug, fmt.Sprintf("%s = %v", statement.PropertyPath, statement.Value))
			}
			builders[i].AddToVeneerTrail(fmt.Sprintf("Initialize[%s]", strings.Join(veneerDebug, ", ")))
		}

		return builders, nil
	}
}

// PromoteOptionsToConstructor promotes the given options as constructor
// parameters. Both arguments and assignments described by the options
// will be exposed in the builder's constructor.
func PromoteOptionsToConstructor(selector Selector, optionNames []string) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			for _, optName := range optionNames {
				opt, ok := builder.OptionByName(optName)
				if !ok {
					continue
				}

				// TODO: do it for every argument/assignment?
				arg := opt.Args[0].DeepCopy()
				arg.Type.Nullable = false

				builders[i].Constructor.Args = append(builders[i].Constructor.Args, arg)
				builders[i].Constructor.Assignments = append(builders[i].Constructor.Assignments, opt.Assignments[0])

				builders[i].AddToVeneerTrail(fmt.Sprintf("PromoteOptionsToConstructor[%s]", optName))
			}
		}

		return builders, nil
	}
}

// AddOption adds a completely new option to the selected builders.
func AddOption(selector Selector, newOption veneers.Option) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			newOpt, err := newOption.AsIR(builders, builder)
			if err != nil {
				return nil, err
			}

			builders[i].Options = append(builders[i].Options, newOpt)
			builders[i].AddToVeneerTrail("AddOption")
		}

		return builders, nil
	}
}

// DefaultToConstant sets a default value into a constant.
// When we unfold a value, omit an option, or any other action; we can lose the default value in the Builder struct.
// If we parse the defaults using the Builder and not the Schema, we might not have all the defaults and with this rule,
// we set these defaults as constants inside the constructor.
func DefaultToConstant(selector Selector, options []string) RewriteRule {
	return func(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
		for i, builder := range builders {
			if !selector(schemas, builder) {
				continue
			}

			for _, option := range options {
				opt, ok := builder.OptionByName(option)
				if !ok {
					continue
				}

				values := make([]any, 0)
				for _, args := range opt.Args {
					if args.Type.Default != nil {
						values = append(values, args.Type.Default)
					}
				}

				for v := range opt.Assignments {
					opt.Assignments[v].Value.Argument = nil
					opt.Assignments[v].Value.Constant = values[v]
				}

				builders[i].Constructor.Assignments = append(builders[i].Constructor.Assignments, opt.Assignments...)
			}

		}

		return builders, nil
	}
}

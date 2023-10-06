package builder

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type RewriteRule func(builders ast.Builders) (ast.Builders, error)

func mapToSelected(selector Selector, mapFunc func(builders ast.Builders, builder ast.Builder) (ast.Builder, error)) RewriteRule {
	return func(builders ast.Builders) (ast.Builders, error) {
		for i, b := range builders {
			if !selector(b) {
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

func mergeOptions(fromBuilder ast.Builder, intoBuilder ast.Builder, underPath ast.Path, excludeOptions []string) (ast.Builder, error) {
	newBuilder := intoBuilder

	for _, opt := range fromBuilder.Options {
		if tools.ItemInList(opt.Name, excludeOptions) {
			continue
		}

		newOpt := opt
		newOpt.Assignments = nil

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
	return func(builders ast.Builders) (ast.Builders, error) {
		filteredBuilders := make([]ast.Builder, 0, len(builders))

		for _, builder := range builders {
			if selector(builder) {
				continue
			}

			filteredBuilders = append(filteredBuilders, builder)
		}

		return filteredBuilders, nil
	}
}

func MergeInto(selector Selector, sourceBuilderName string, underPath string, excludeOptions []string) RewriteRule {
	return mapToSelected(selector, func(builders ast.Builders, destinationBuilder ast.Builder) (ast.Builder, error) {
		sourcePkg, sourceBuilderNameWithoutPkg, found := strings.Cut(sourceBuilderName, ".")
		if !found {
			return destinationBuilder, fmt.Errorf("sourceBuilderName '%s' is incorrect: no package found", sourceBuilderName)
		}

		sourceBuilder, found := builders.LocateByObject(sourcePkg, sourceBuilderNameWithoutPkg)
		if !found {
			return destinationBuilder, fmt.Errorf("source builder '%s.%s' not found", sourcePkg, sourceBuilderNameWithoutPkg)
		}

		newRoot, err := destinationBuilder.MakePath(builders, underPath)
		if err != nil {
			return destinationBuilder, err
		}

		// TODO: initializations
		return mergeOptions(sourceBuilder, destinationBuilder, newRoot, excludeOptions)
	})
}

func composePanelType(builders ast.Builders, panelType string, panelBuilder ast.Builder, composableBuilders ast.Builders) (ast.Builder, error) {
	newBuilder := ast.Builder{
		Schema:      panelBuilder.Schema,
		For:         panelBuilder.For,
		RootPackage: panelType,
		Package:     "panel",
	}

	typeField, ok := panelBuilder.For.Type.AsStruct().FieldByName("type")
	if !ok {
		return panelBuilder, fmt.Errorf("could not find field 'type' in panel builder")
	}

	newBuilder.Initializations = append(newBuilder.Initializations, ast.Assignment{
		Path: ast.Path{
			{
				Identifier: "type",
				Type:       typeField.Type,
			},
		},
		Value: panelType,
	})

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

		newBuilder.Options = append(newBuilder.Options, panelOpt)
	}

	compositionMap := map[string]string{ // Builder → assignment root path
		"Options":     "options",
		"FieldConfig": "fieldConfig.defaults.custom",
	}

	for _, composableBuilder := range composableBuilders {
		underPath, exists := compositionMap[composableBuilder.For.Name]
		if !exists {
			return newBuilder, fmt.Errorf("unexpected composable type '%s'", composableBuilder.For.Name)
		}

		newRoot, err := newBuilder.MakePath(builders, underPath)
		if err != nil {
			return newBuilder, err
		}

		newBuilder.Initializations = append(newBuilder.Initializations, ast.Assignment{
			Path:       newRoot,
			EmptyValue: true,
		})

		ref := composableBuilder.For.SelfRef
		refType := ast.NewRef(ref.ReferredPkg, ref.ReferredType)
		newRoot[len(newRoot)-1].TypeHint = &refType

		newBuilder, err = mergeOptions(composableBuilder, newBuilder, newRoot, nil)
		if err != nil {
			return newBuilder, err
		}
	}

	return newBuilder, nil
}

func ComposeDashboardPanel(selector Selector, panelBuilderName string) RewriteRule {
	return func(builders ast.Builders) (ast.Builders, error) {
		panelBuilderPkg, panelBuilderNameWithoutPkg, found := strings.Cut(panelBuilderName, ".")
		if !found {
			return nil, fmt.Errorf("panelBuilderName '%s' is incorrect: no package found", panelBuilderPkg)
		}

		panelBuilder, found := builders.LocateByObject(panelBuilderPkg, panelBuilderNameWithoutPkg)
		if !found {
			return nil, fmt.Errorf("panel builder '%s' not found", panelBuilderName)
		}

		// - add to newBuilders all the builders that are not composable (ie: don't comply to the selector)
		// - build a map of composable builders, indexed by panel type
		// - aggregate the composable builders into a new, composed panel builder
		// - add the new composed panel builders to newBuilders

		newBuilders := make([]ast.Builder, 0, len(builders))
		composableBuilders := make(map[string]ast.Builders)

		for _, builder := range builders {
			// the builder is for a composable type
			if selector(builder) {
				panelType := builder.Schema.Metadata.Identifier
				composableBuilders[panelType] = append(composableBuilders[panelType], builder)
				continue
			}

			newBuilders = append(newBuilders, builder)
		}

		for panelType, buildersForType := range composableBuilders {
			composedBuilder, err := composePanelType(builders, panelType, panelBuilder, buildersForType)
			if err != nil {
				return nil, err
			}

			newBuilders = append(newBuilders, composedBuilder)
		}

		return newBuilders, nil
	}
}

package builder

import (
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
)

type RewriteAction func(builders ast.Builders, builder ast.Builder) ast.Builder

func OmitAction() RewriteAction {
	return func(builders ast.Builders, _ ast.Builder) ast.Builder {
		return ast.Builder{}
	}
}

func MergeIntoAction(sourceBuilderName string, underPath string, excludeOptions []string) RewriteAction {
	return func(builders ast.Builders, destinationBuilder ast.Builder) ast.Builder {
		sourcePkg, sourceBuilderNameWithoutPkg, found := strings.Cut(sourceBuilderName, ".")
		if !found {
			sourcePkg = destinationBuilder.Package
			sourceBuilderNameWithoutPkg = sourceBuilderName
		}

		sourceBuilder, found := builders.LocateByObject(sourcePkg, sourceBuilderNameWithoutPkg)
		if !found {
			return destinationBuilder
		}

		newBuilder := destinationBuilder

		// TODO: initializations

		for _, opt := range sourceBuilder.Options {
			if tools.ItemInList(opt.Name, excludeOptions) {
				continue
			}

			// TODO: assignment paths
			newOpt := opt
			newOpt.Assignments = nil

			for _, assignment := range opt.Assignments {
				newAssignment := assignment
				// @FIXME: this only works if no part of the `underPath` path can be nil
				newAssignment.Path = underPath + "." + assignment.Path

				newOpt.Assignments = append(newOpt.Assignments, newAssignment)
			}

			newBuilder.Options = append(newBuilder.Options, newOpt)
		}

		return newBuilder
	}
}

func ComposeDashboardPanelAction(panelBuilderName string) RewriteAction {
	return func(builders ast.Builders, destinationBuilder ast.Builder) ast.Builder {
		panelBuilderPkg, panelBuilderNameWithoutPkg, found := strings.Cut(panelBuilderName, ".")
		if !found {
			panelBuilderPkg = destinationBuilder.Package
			panelBuilderNameWithoutPkg = panelBuilderName
		}

		panelBuilder, found := builders.LocateByObject(panelBuilderPkg, panelBuilderNameWithoutPkg)
		if !found {
			// TODO: we failed here, an error is the correct action.
			return destinationBuilder
		}

		// we're more building a new thing than updating any of the panel or composable builder,
		// so let's create a completely fresh instance
		newBuilder := ast.Builder{
			Schema:      panelBuilder.Schema,
			For:         panelBuilder.For,
			RootPackage: destinationBuilder.RootPackage,
			Package:     destinationBuilder.Package,
		}

		// re-add panel-related options
		for _, panelOpt := range panelBuilder.Options {
			// this value is a constant that depends on the plugin being composed into a panel
			if panelOpt.Name == "type" {
				continue
			}

			newBuilder.Options = append(newBuilder.Options, panelOpt)
		}

		// re-add plugin-specific options
		for _, pluginOpt := range destinationBuilder.Options {
			newOpt := pluginOpt
			newOpt.Assignments = nil

			// re-root assignments
			for _, assignment := range pluginOpt.Assignments {
				newAssignment := assignment
				newAssignment.Path = "options." + newAssignment.Path

				newOpt.Assignments = append(newOpt.Assignments, newAssignment)
			}

			newBuilder.Options = append(newBuilder.Options, newOpt)
		}

		newBuilder.Initializations = append(newBuilder.Initializations, ast.Assignment{
			Path:      "type",
			ValueType: ast.String(),
			Value:     destinationBuilder.Schema.Metadata.Identifier,
		})

		return newBuilder
	}
}

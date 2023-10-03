package builder

import (
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
		// we're implicitly saying that this action only works on builders originating from the same package.
		// that's probably not good enough.

		sourceBuilder, found := builders.LocateByObject(destinationBuilder.Package, sourceBuilderName)
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

func ComposeDashboardPanelAction(_ /*panelBuilderName*/ string) RewriteAction {
	return func(builders ast.Builders, destinationBuilder ast.Builder) ast.Builder {
		return destinationBuilder
	}
}

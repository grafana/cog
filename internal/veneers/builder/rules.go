package builder

import (
	"fmt"
	"strings"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers"
)

type Rule struct {
	Selector *Selector
	Action   *Action
}

func (rule Rule) Matches(schemas ast.Schemas, builder ast.Builder) bool {
	return rule.Selector.Matches(schemas, builder)
}

func (rule Rule) Apply(ctx RuleCtx, selectedBuilders ast.Builders) (ast.Builders, error) {
	return rule.Action.run(ctx, selectedBuilders)
}

func (rule Rule) String() string {
	return fmt.Sprintf("selector=%s, action=%s", rule.Selector, rule.Action)
}

type RuleCtx struct {
	Schemas  ast.Schemas
	Builders ast.Builders
}

type ActionRunner func(ctx RuleCtx, selectedBuilders ast.Builders) (ast.Builders, error)

type Action struct {
	description string
	run         ActionRunner
}

func (action Action) Run(ctx RuleCtx, selectedBuilders ast.Builders) (ast.Builders, error) {
	return action.run(ctx, selectedBuilders)
}

func (action Action) String() string {
	return action.description
}

// Omit removes a builder.
func Omit(selector *Selector) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: "omit",
			run: func(ctx RuleCtx, builders ast.Builders) (ast.Builders, error) {
				return nil, nil
			},
		},
	}
}

func MergeInto(selector *Selector, sourceBuilderName string, underPath string, excludeOptions []string, renameOptions map[string]string) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("merge_into[source: '%s', under_path: '%s']", sourceBuilderName, underPath),
			run: MergeIntoAction(
				sourceBuilderName,
				underPath,
				excludeOptions,
				renameOptions,
			),
		},
	}
}

func ComposeBuilders(selector *Selector, config CompositionConfig) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("compose_builders[source: '%s']", config.SourceBuilderName),
			run:         ComposeBuildersAction(config),
		},
	}
}

// Rename renames a builder.
func Rename(selector *Selector, newName string) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("rename[as: '%s']", newName),
			run:         RenameAction(newName),
		},
	}
}

func VeneerTrailAsComments(selector *Selector) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: "venner_trail_as_comments",
			run:         VeneerTrailAsCommentsAction(),
		},
	}
}

func Properties(selector *Selector, properties []ast.StructField) *Rule {
	propNames := tools.Map(properties, func(prop ast.StructField) string {
		return prop.Name
	})

	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("properties['%s']", strings.Join(propNames, ", ")),
			run:         PropertiesAction(properties),
		},
	}
}

// Duplicate duplicates a builder.
// The name of the duplicated builder has to be specified and some options can
// be excluded.
func Duplicate(selector *Selector, duplicateName string, excludeOptions []string) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("duplicate[as: '%s']", duplicateName),
			run:         DuplicateAction(duplicateName, excludeOptions),
		},
	}
}

func Initialize(selector *Selector, statements []Initialization) *Rule {
	initPaths := tools.Map(statements, func(stmt Initialization) string {
		return stmt.PropertyPath
	})

	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("initialize[%s]", strings.Join(initPaths, ", ")),
			run:         InitializeAction(statements),
		},
	}
}

// PromoteOptionsToConstructor promotes the given options as constructor
// parameters. Both arguments and assignments described by the options
// will be exposed in the builder's constructor.
func PromoteOptionsToConstructor(selector *Selector, optionNames []string) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("promote_options_to_constructor[opts: '%s']", strings.Join(optionNames, ", ")),
			run:         PromoteOptionsToConstructorAction(optionNames),
		},
	}
}

// AddOption adds a completely new option to the selected builders.
func AddOption(selector *Selector, newOption veneers.Option) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("add_option[name: '%s']", newOption.Name),
			run:         AddOptionAction(newOption),
		},
	}
}

// AddFactory adds a builder factory to the selected builders.
// These factories are meant to be used to simplify the instantiation of
// builders for common use-cases.
func AddFactory(selector *Selector, factory ast.BuilderFactory) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: fmt.Sprintf("add_factory[name: '%s']", factory.Name),
			run:         AddFactoryAction(factory),
		},
	}
}

// Debug prints debugging information about a builder.
func Debug(selector *Selector) *Rule {
	return &Rule{
		Selector: selector,
		Action: &Action{
			description: "debug",
			run:         DebugAction(),
		},
	}
}

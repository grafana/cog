package yaml

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/tools"
	"github.com/grafana/cog/internal/veneers/builder"
)

/******************************************************************************
 * Rules
 *****************************************************************************/

type BuilderRule struct {
	Omit                  *BuilderSelector       `yaml:"omit"`
	Rename                *RenameBuilder         `yaml:"rename"`
	MergeInto             *MergeInto             `yaml:"merge_into"`
	ComposeDashboardPanel *ComposeDashboardPanel `yaml:"compose_dashboard_panel"`
	Properties            *Properties            `yaml:"properties"`
	Duplicate             *Duplicate             `yaml:"duplicate"`
	Initialize            *Initialize            `yaml:"initialize"`
}

func (rule BuilderRule) AsRewriteRule(pkg string) (builder.RewriteRule, error) {
	if rule.Omit != nil {
		selector, err := rule.Omit.AsSelector(pkg)
		if err != nil {
			return nil, err
		}

		return builder.Omit(selector), nil
	}

	if rule.Rename != nil {
		return rule.Rename.AsRewriteRule(pkg)
	}

	if rule.MergeInto != nil {
		return rule.MergeInto.AsRewriteRule(pkg)
	}

	if rule.ComposeDashboardPanel != nil {
		return rule.ComposeDashboardPanel.AsRewriteRule()
	}

	if rule.Properties != nil {
		return rule.Properties.AsRewriteRule(pkg)
	}

	if rule.Duplicate != nil {
		return rule.Duplicate.AsRewriteRule(pkg)
	}

	if rule.Initialize != nil {
		return rule.Initialize.AsRewriteRule(pkg)
	}

	return nil, fmt.Errorf("empty rule")
}

type RenameBuilder struct {
	BuilderSelector `yaml:",inline"`

	As string `yaml:"as"`
}

func (rule RenameBuilder) AsRewriteRule(pkg string) (builder.RewriteRule, error) {
	selector, err := rule.AsSelector(pkg)
	if err != nil {
		return nil, err
	}

	return builder.Rename(selector, rule.As), nil
}

type MergeInto struct {
	Destination    string   `yaml:"destination"`
	Source         string   `yaml:"source"`
	UnderPath      string   `yaml:"under_path"`
	ExcludeOptions []string `yaml:"exclude_options"`
}

func (rule MergeInto) AsRewriteRule(pkg string) (builder.RewriteRule, error) {
	return builder.MergeInto(
		builder.ByObjectName(pkg, rule.Destination),
		rule.Source,
		rule.UnderPath,
		rule.ExcludeOptions,
	), nil
}

type ComposeDashboardPanel struct {
	PanelBuilderName    string   `yaml:"panel_builder_name"`
	ExcludePanelOptions []string `yaml:"exclude_panel_options"`
}

func (rule ComposeDashboardPanel) AsRewriteRule() (builder.RewriteRule, error) {
	return builder.ComposeDashboardPanel(
		builder.ComposableDashboardPanel(),
		rule.PanelBuilderName,
		rule.ExcludePanelOptions,
	), nil
}

type Properties struct {
	BuilderSelector `yaml:",inline"`
	Set             []ast.StructField `yaml:"set"`
}

func (rule Properties) AsRewriteRule(pkg string) (builder.RewriteRule, error) {
	selector, err := rule.AsSelector(pkg)
	if err != nil {
		return nil, err
	}

	return builder.Properties(
		selector,
		rule.Set,
	), nil
}

type Duplicate struct {
	BuilderSelector `yaml:",inline"`
	As              string   `yaml:"as"`
	ExcludeOptions  []string `yaml:"exclude_options"`
}

func (rule Duplicate) AsRewriteRule(pkg string) (builder.RewriteRule, error) {
	selector, err := rule.AsSelector(pkg)
	if err != nil {
		return nil, err
	}

	return builder.Duplicate(
		selector,
		rule.As,
		rule.ExcludeOptions,
	), nil
}

type Initialization struct {
	Property string `yaml:"property"`
	Value    any    `yaml:"value"`
}

type Initialize struct {
	BuilderSelector `yaml:",inline"`
	Set             []Initialization `yaml:"set"`
}

func (rule Initialize) AsRewriteRule(pkg string) (builder.RewriteRule, error) {
	selector, err := rule.AsSelector(pkg)
	if err != nil {
		return nil, err
	}

	return builder.Initialize(
		selector,
		tools.Map(rule.Set, func(init Initialization) builder.Initialization {
			return builder.Initialization{PropertyPath: init.Property, Value: init.Value}
		}),
	), nil
}

/******************************************************************************
 * Selectors
 *****************************************************************************/

type BuilderSelector struct {
	ByObject *string `yaml:"by_object"`
	ByName   *string `yaml:"by_name"`

	GeneratedFromDisjunction *bool `yaml:"generated_from_disjunction"` // noop?
}

func (selector BuilderSelector) AsSelector(pkg string) (builder.Selector, error) {
	if selector.ByObject != nil {
		return builder.ByObjectName(pkg, *selector.ByObject), nil
	}

	if selector.ByName != nil {
		return builder.ByName(pkg, *selector.ByName), nil
	}

	if selector.GeneratedFromDisjunction != nil {
		return builder.StructGeneratedFromDisjunction(), nil
	}

	return nil, fmt.Errorf("empty selector")
}

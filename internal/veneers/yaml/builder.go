package yaml

import (
	"fmt"

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
}

func (rule BuilderRule) AsRewriteRule() (builder.RewriteRule, error) {
	if rule.Omit != nil {
		selector, err := rule.Omit.AsSelector()
		if err != nil {
			return nil, err
		}

		return builder.Omit(selector), nil
	}

	if rule.Rename != nil {
		return rule.Rename.AsRewriteRule()
	}

	if rule.MergeInto != nil {
		return rule.MergeInto.AsRewriteRule()
	}

	if rule.ComposeDashboardPanel != nil {
		return rule.ComposeDashboardPanel.AsRewriteRule()
	}

	return nil, fmt.Errorf("empty rule")
}

type RenameBuilder struct {
	BuilderSelector `yaml:",inline"`

	As string `yaml:"as"`
}

func (rule RenameBuilder) AsRewriteRule() (builder.RewriteRule, error) {
	selector, err := rule.AsSelector()
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

func (rule MergeInto) AsRewriteRule() (builder.RewriteRule, error) {
	return builder.MergeInto(
		builder.ByObjectName(rule.Destination),
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

/******************************************************************************
 * Selectors
 *****************************************************************************/

type BuilderSelector struct {
	ByObject *string `yaml:"by_object"`

	GeneratedFromDisjunction *bool `yaml:"generated_from_disjunction"` // noop?
	Any                      *bool `yaml:"any"`                        // noop?
}

func (selector BuilderSelector) AsSelector() (builder.Selector, error) {
	if selector.ByObject != nil {
		return builder.ByObjectName(*selector.ByObject), nil
	}

	if selector.GeneratedFromDisjunction != nil {
		return builder.StructGeneratedFromDisjunction(), nil
	}

	if selector.Any != nil {
		return builder.EveryBuilder(), nil
	}

	return nil, fmt.Errorf("empty selector")
}

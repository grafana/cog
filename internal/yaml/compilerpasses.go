package yaml

import (
	"fmt"

	"github.com/grafana/cog/internal/ast/compiler"
)

type CompilerPass struct {
	DataqueryIdentification *DataqueryIdentification `yaml:"dataquery_identification"`
	Unspec                  *Unspec                  `yaml:"unspec"`

	Dashboard           *Dashboard           `yaml:"dashboard"`
	DashboardPanels     *DashboardPanels     `yaml:"dashboard_panels"`
	DashboardTargets    *DashboardTargets    `yaml:"dashboard_targets"`
	DashboardTimePicker *DashboardTimePicker `yaml:"dashboard_timepicker"`

	Cloudwatch            *Cloudwatch            `yaml:"cloudwatch"`
	GoogleCloudMonitoring *GoogleCloudMonitoring `yaml:"google_cloud_monitoring"`
	PrometheusDataquery   *PrometheusDataquery   `yaml:"prometheus_dataquery"`
	LibraryPanels         *LibraryPanels         `yaml:"library_panels"`
	TestData              *TestData              `yaml:"test_data"`
}

func (pass CompilerPass) AsCompilerPass() (compiler.Pass, error) {
	if pass.DataqueryIdentification != nil {
		return pass.DataqueryIdentification.AsCompilerPass(), nil
	}
	if pass.Unspec != nil {
		return pass.Unspec.AsCompilerPass(), nil
	}

	if pass.Dashboard != nil {
		return pass.Dashboard.AsCompilerPass(), nil
	}
	if pass.DashboardPanels != nil {
		return pass.DashboardPanels.AsCompilerPass(), nil
	}
	if pass.DashboardTargets != nil {
		return pass.DashboardTargets.AsCompilerPass(), nil
	}
	if pass.DashboardTimePicker != nil {
		return pass.DashboardTimePicker.AsCompilerPass(), nil
	}

	if pass.Cloudwatch != nil {
		return pass.Cloudwatch.AsCompilerPass(), nil
	}
	if pass.GoogleCloudMonitoring != nil {
		return pass.GoogleCloudMonitoring.AsCompilerPass(), nil
	}
	if pass.PrometheusDataquery != nil {
		return pass.PrometheusDataquery.AsCompilerPass(), nil
	}
	if pass.LibraryPanels != nil {
		return pass.LibraryPanels.AsCompilerPass(), nil
	}
	if pass.TestData != nil {
		return pass.TestData.AsCompilerPass(), nil
	}

	return nil, fmt.Errorf("empty compiler pass")
}

type DataqueryIdentification struct {
}

func (pass DataqueryIdentification) AsCompilerPass() compiler.Pass {
	return &compiler.DataqueryIdentification{}
}

type Unspec struct {
}

func (pass Unspec) AsCompilerPass() compiler.Pass {
	return &compiler.Unspec{}
}

type Dashboard struct {
}

func (pass Dashboard) AsCompilerPass() compiler.Pass {
	return &compiler.Dashboard{}
}

type DashboardPanels struct {
}

func (pass DashboardPanels) AsCompilerPass() compiler.Pass {
	return &compiler.DashboardPanelsRewrite{}
}

type DashboardTargets struct {
}

func (pass DashboardTargets) AsCompilerPass() compiler.Pass {
	return &compiler.DashboardTargetsRewrite{}
}

type DashboardTimePicker struct {
}

func (pass DashboardTimePicker) AsCompilerPass() compiler.Pass {
	return &compiler.DashboardTimePicker{}
}

type Cloudwatch struct {
}

func (pass Cloudwatch) AsCompilerPass() compiler.Pass {
	return &compiler.Cloudwatch{}
}

type GoogleCloudMonitoring struct {
}

func (pass GoogleCloudMonitoring) AsCompilerPass() compiler.Pass {
	return &compiler.GoogleCloudMonitoring{}
}

type PrometheusDataquery struct {
}

func (pass PrometheusDataquery) AsCompilerPass() compiler.Pass {
	return &compiler.PrometheusDataquery{}
}

type LibraryPanels struct {
}

func (pass LibraryPanels) AsCompilerPass() compiler.Pass {
	return &compiler.LibraryPanels{}
}

type TestData struct {
}

func (pass TestData) AsCompilerPass() compiler.Pass {
	return &compiler.TestData{}
}

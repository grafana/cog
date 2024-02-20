package yaml

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
)

type CompilerPass struct {
	DataqueryIdentification *DataqueryIdentification `yaml:"dataquery_identification"`
	Unspec                  *Unspec                  `yaml:"unspec"`
	FieldsSetRequired       *FieldsSetRequired       `yaml:"fields_set_required"`
	Omit                    *Omit                    `yaml:"omit"`
	AddFields               *AddFields               `yaml:"add_fields"`

	DashboardPanels     *DashboardPanels     `yaml:"dashboard_panels"`
	DashboardTargets    *DashboardTargets    `yaml:"dashboard_targets"`
	DashboardTimePicker *DashboardTimePicker `yaml:"dashboard_timepicker"`

	Cloudwatch            *Cloudwatch            `yaml:"cloudwatch"`
	GoogleCloudMonitoring *GoogleCloudMonitoring `yaml:"google_cloud_monitoring"`
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

	if pass.FieldsSetRequired != nil {
		return pass.FieldsSetRequired.AsCompilerPass()
	}
	if pass.Omit != nil {
		return pass.Omit.AsCompilerPass()
	}
	if pass.AddFields != nil {
		return pass.AddFields.AsCompilerPass()
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

type FieldsSetRequired struct {
	Fields []string // Expected format: [package].[object].[field]
}

func (pass FieldsSetRequired) AsCompilerPass() (compiler.Pass, error) {
	fieldRefs := make([]compiler.FieldReference, 0, len(pass.Fields))

	for _, ref := range pass.Fields {
		fieldRef, err := compiler.FieldReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		fieldRefs = append(fieldRefs, fieldRef)
	}

	return &compiler.FieldsSetRequired{Fields: fieldRefs}, nil
}

type Omit struct {
	Objects []string // Expected format: [package].[object]
}

func (pass Omit) AsCompilerPass() (compiler.Pass, error) {
	objectRefs := make([]compiler.ObjectReference, 0, len(pass.Objects))

	for _, ref := range pass.Objects {
		objectRef, err := compiler.ObjectReferenceFromString(ref)
		if err != nil {
			return nil, err
		}

		objectRefs = append(objectRefs, objectRef)
	}

	return &compiler.Omit{Objects: objectRefs}, nil
}

type AddFields struct {
	To     string // Expected format: [package].[object]
	Fields []ast.StructField
}

func (pass AddFields) AsCompilerPass() (compiler.Pass, error) {
	objectRef, err := compiler.ObjectReferenceFromString(pass.To)
	if err != nil {
		return nil, err
	}

	return &compiler.AddFields{
		Object: objectRef,
		Fields: pass.Fields,
	}, nil
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

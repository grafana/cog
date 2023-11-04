package golang

import (
	"github.com/grafana/codejen"
	"github.com/grafana/cog/internal/jennies/context"
)

type BuilderInterface struct {
}

func (jenny BuilderInterface) JennyName() string {
	return "GolangOptionsBuilder"
}

func (jenny BuilderInterface) Generate(_ context.Builders) (*codejen.File, error) {
	output := jenny.generateFile()

	return codejen.NewFile("builder.go", []byte(output), jenny), nil
}

func (jenny BuilderInterface) generateFile() string {
	return `package cog

import (
	"fmt"
)

type BuildErrors []*BuildError

func (errs BuildErrors) Error() string {
	var b []byte
	for i, err := range errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

type BuildError struct {
	Path string
	Message string
}

func (err *BuildError) Error() string {
	return fmt.Sprintf("%s: %s", err.Path, err.Message)
}

func MakeBuildErrors(rootPath string, err error) BuildErrors {
	if buildErrs, ok := err.(BuildErrors); ok {
		for _, buildErr := range buildErrs {
			buildErr.Path = rootPath + "." + buildErr.Path
		}

		return buildErrs
	}
	
	if buildErr, ok := err.(*BuildError); ok {
		return BuildErrors{buildErr}
	}

	return BuildErrors{&BuildError{
		Path:    rootPath,
		Message: err.Error(),
	}}
}

type Builder[ResourceT any] interface {
  Build() (ResourceT, error)
}

type PanelcfgVariantMarshalConfig struct {
	OptionsUnmarshaler     func(raw []byte) (any, error)
	FieldConfigUnmarshaler func(raw []byte) (any, error)
}

var panelcfgVariantsRegistry map[string]PanelcfgVariantMarshalConfig

func RegisterPanelcfgVariant(variantName string, config PanelcfgVariantMarshalConfig) {
	if panelcfgVariantsRegistry == nil {
		panelcfgVariantsRegistry = make(map[string]PanelcfgVariantMarshalConfig)
	}
	panelcfgVariantsRegistry[variantName] = config
}

func MarshalConfigForPanelcfgVariant(variantName string) (PanelcfgVariantMarshalConfig, bool) {
	config, found := panelcfgVariantsRegistry[variantName]

	return config, found
}

type DataqueryVariantMarshalConfig struct {
	DataqueryUnmarshaler func(raw []byte) (interface {
	ImplementsDataqueryVariant()
}, error)
}

var dataqueryVariantsRegistry []DataqueryVariantMarshalConfig

func RegisterDataqueryVariant(config DataqueryVariantMarshalConfig) {
	dataqueryVariantsRegistry = append(dataqueryVariantsRegistry, config)
}

func UnmarshalDataqueryArray[VariantT interface {
	ImplementsDataqueryVariant()
}](raw []byte) ([]VariantT, error) {
	rawDataqueries := []json.RawMessage{}
	if err := json.Unmarshal(raw, &rawDataqueries); err != nil {
		return nil, err
	}

	dataqueries := make([]VariantT, 0, len(rawDataqueries))
	for _, rawDataquery := range rawDataqueries {
		dataquery, err := UnmarshalDataquery[VariantT](rawDataquery)
		if err != nil {
			return nil, err
		}

		dataqueries = append(dataqueries, dataquery)
	}

	return dataqueries, nil
}

func UnmarshalDataquery[VariantT any](raw []byte) (VariantT, error) {
	var empty VariantT

	for _, config := range dataqueryVariantsRegistry {
		if config.DataqueryUnmarshaler == nil {
			continue
		}

		dataquery, err := config.DataqueryUnmarshaler(raw)
		if err != nil {
			continue
		}

		return dataquery.(VariantT), nil
	}

	return empty, fmt.Errorf("could not unmarshal dataquery")
}

`
}

package yaml

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type ConverterConfig struct {
	Runtime []Runtime `yaml:"runtime"`
}

type Runtime struct {
	Package            string `yaml:"package"`
	Name               string `yaml:"name"`
	NameFunc           string `yaml:"name_func"`
	DiscriminatorField string `yaml:"discriminator_field"`
}

type ConverterConfigReader struct{}

func NewConverterConfigReader() *ConverterConfigReader {
	return &ConverterConfigReader{}
}

func (c ConverterConfigReader) ReadConverterConfig(filename string) (ConverterConfig, error) {
	if filename == "" {
		return ConverterConfig{}, nil
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		return ConverterConfig{}, err
	}

	var config ConverterConfig
	if err = yaml.UnmarshalWithOptions(contents, &config, yaml.DisallowUnknownField()); err != nil {
		return config, fmt.Errorf("can not read converter config: %s\n%s", filename, yaml.FormatError(err, true, true))
	}

	return config, nil
}

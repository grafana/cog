package codegen

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/grafana/cog/internal/tools"
)

// Unit represents part of a codegen pipeline.
// A unit describes a set of inputs and their associated schema/builder
// transformations.
type Unit struct {
	Inputs []*Input `yaml:"inputs"`

	// BuilderTransforms holds a list of paths to builder transformation files.
	// The paths can refer to files or directories.
	BuilderTransforms []string `yaml:"builder_transformations"`
}

// unitFromFile loads a codegen unit from the given file.
func unitFromFile(file string, parameters map[string]string) (*Unit, error) {
	var err error
	if !filepath.IsAbs(file) {
		file, err = filepath.Abs(file)
		if err != nil {
			return nil, err
		}
	}

	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	unit := &Unit{}
	if err := yaml.UnmarshalWithOptions(contents, unit, yaml.DisallowUnknownField()); err != nil {
		return nil, fmt.Errorf("can not load codegen unit: %s\n%s", file, yaml.FormatError(err, true, true))
	}

	unitParams := map[string]string{}
	maps.Copy(unitParams, parameters)
	unitParams["__config_dir"] = filepath.Dir(file)

	unit.interpolateParameters(createInterpolator(unitParams))

	return unit, nil
}

func (unit *Unit) interpolateParameters(interpolator ParametersInterpolator) {
	for _, input := range unit.Inputs {
		// An error can only happen with the input isn't descriptive.
		// This case should have already been handled before
		// interpolateParameters() is called.
		_ = input.InterpolateParameters(interpolator)
	}

	unit.BuilderTransforms = tools.Map(unit.BuilderTransforms, interpolator)
}

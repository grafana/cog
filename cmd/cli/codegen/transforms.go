package codegen

import (
	"github.com/grafana/cog/internal/tools"
)

type Transforms struct {
	// CompilerPassesFiles holds a list of paths to files containing compiler
	// passes to apply to all the schemas.
	CompilerPassesFiles []string `yaml:"schemas"`

	// VeneersDirectories holds a list of paths to directories containing
	// veneers to apply to all the builders.
	VeneersDirectories []string `yaml:"builders"`
}

func (transforms *Transforms) interpolateParameters(interpolator ParametersInterpolator) {
	transforms.CompilerPassesFiles = tools.Map(transforms.CompilerPassesFiles, interpolator)
	transforms.VeneersDirectories = tools.Map(transforms.VeneersDirectories, interpolator)
}

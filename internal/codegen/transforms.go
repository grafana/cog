package codegen

import (
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/tools"
)

type VeneerPath struct {
	If   string `yaml:"if"`
	Path string `yaml:"path"`
}

type Transforms struct {
	// CommonPassesFiles holds a list of paths to files containing compiler
	// passes to apply to all the schemas.
	// Note: these compiler passes are applied *before* language-specific passes.
	CommonPassesFiles []string `yaml:"schemas"`

	// CommonPasses holds a list of compiler passes to apply to all the schemas.
	// If this field is set, CommonPassesFiles is ignored.
	// Note: these compiler passes are applied *before* language-specific passes.
	CommonPasses compiler.Passes `yaml:"-"`

	// FinalPasses holds a list of compiler passes to apply to all the schemas.
	// Note: these compiler passes are applied *after* language-specific passes.
	FinalPasses compiler.Passes `yaml:"-"`

	// VeneersDirectories holds a list of paths to directories containing
	// veneers to apply to all the builders.
	VeneersDirectories []VeneerPath `yaml:"builders"`

	// ConverterConfig is the configuration modify the converters output.
	ConverterConfig string `yaml:"converters"`
}

func (transforms *Transforms) interpolateParameters(interpolator ParametersInterpolator) {
	transforms.CommonPassesFiles = tools.Map(transforms.CommonPassesFiles, interpolator)
	transforms.VeneersDirectories = tools.Map(transforms.VeneersDirectories, func(path VeneerPath) VeneerPath {
		path.If = interpolator(path.If)
		path.Path = interpolator(path.Path)

		return path
	})
	transforms.ConverterConfig = interpolator(transforms.ConverterConfig)
}

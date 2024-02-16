package common

// Config represents global configuration options, to be used by all jennies.
type Config struct {
	// Debug turns on or off the debugging mode.
	Debug bool

	// Types indicates whether types should be generated or not.
	Types bool

	// Builders indicates whether builders should be generated or not.
	Builders bool

	// Optional filename output name
	OutputFilename string

	// Optional path for runtime, variants, etc. necessary code
	UtilsPath string

	// Path for package imports
	PackageRoot string
}

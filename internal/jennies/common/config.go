package common

// Config represents global configuration options, to be used by all jennies.
type Config struct {
	// Debug turns on or off the debugging mode.
	Debug bool

	// Types indicates whether types should be generated or not.
	Types bool

	// Builders indicates whether builders should be generated or not.
	Builders bool

	RenameOutputFunc func(pkg string) string

	GoConfig GoConfig
	TSConfig TSConfig
}

type GoConfig struct {
	// PackageRoot defines the package imports for Go
	PackageRoot string
}

type TSConfig struct {
	// GenTSIndex indicates whether ts index should be generated or not.
	GenTSIndex   bool
	GenRuntime   bool
	RuntimePath  *string
	ImportMapper func(pkg string) string
}

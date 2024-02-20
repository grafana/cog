package cog

type Kind string

const CoreKind Kind = "core"
const ComposableKind Kind = "composable"
const CustomKind Kind = "custom"

type Language string

const GoLanguage Language = "go"
const TSLanguage Language = "typescript"

type Languages []Language

func (l Languages) languages() []string {
	languages := make([]string, len(l))
	for i, language := range l {
		languages[i] = string(language)
	}

	return languages
}

type Config struct {
	Debug            bool
	FileDirs         []string
	OutputDir        string
	Languages        Languages
	Kind             Kind
	RenameOutputFunc func(pkg string) string

	GoConfig GoConfig
	TSConfig TSConfig
}

type GoConfig struct {
	PackageRoot string
}

type TSConfig struct {
	GenTSIndex   bool
	GenRuntime   bool
	RuntimePath  *string
	ImportMapper func(pkg string) string
}

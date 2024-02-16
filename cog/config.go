package cog

type Kind string

const CoreKind Kind = "core"
const ComposableKind Kind = "composable"
const CustomKind Kind = "custom"

type Language string

const GoLanguage Language = "go"
const TSLanguage Language = "ts"

type Languages []Language

func (l Languages) languages() []string {
	languages := make([]string, len(l))
	for i, language := range l {
		languages[i] = string(language)
	}

	return languages
}

type Config struct {
	Debug     bool
	FileDirs  []string
	OutputDir string
	Languages Languages
	Kind      Kind
	// Go configuration
	PackageRoot string
	UtilsPath   string
}

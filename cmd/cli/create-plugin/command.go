package createplugin

import (
	"context"
	"embed"
	"fmt"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/jennies"
	"github.com/grafana/cog/pkg/template"
	"github.com/spf13/cobra"
)

//go:embed templates/*.tmpl
//nolint:gochecknoglobals
var templatesFS embed.FS

type options struct {
	goModule     string
	languageName string
	outputDir    string
	cogVersion   string
}

func Command(version string) *cobra.Command {
	opts := options{
		cogVersion: version,
	}

	cmd := &cobra.Command{
		Use:   "create-plugin LANGUAGE_NAME",
		Short: "Create a new language plugin.",
		Long: `Create a new language plugin.

	# Creates a new language plugin in the './plugin' directory:

	cog create-plugin rust

	# Creates a new language plugin in a specific directory:

	cog create-plugin -o ./rust-plugin rust
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.languageName = args[0]
			return doScaffold(cmd.Context(), opts)
		},
	}

	cmd.Flags().StringVar(&opts.goModule, "go-module-path", "github.com/org/plugin", "Type of intermediate representation to Inspect. Valid values: types, builders, converters.")
	cmd.Flags().StringVarP(&opts.outputDir, "output", "o", "./plugin", "Output directory.")

	return cmd
}

func doScaffold(ctx context.Context, opts options) error {
	tmpl, err := template.New("create-plugin", template.ParseFS(templatesFS, "templates"))
	if err != nil {
		return err
	}

	fs := codejen.NewFS()
	tmplData := map[string]any{
		"LanguageName": opts.languageName,
		"GoModule":     opts.goModule,
		"CogVersion":   opts.cogVersion,
	}
	if opts.cogVersion == "SNAPSHOT" {
		tmplData["CogVersion"] = "latest"
	}

	targets := []struct {
		tmplFile       string
		targetFileName string
	}{
		{tmplFile: "go.mod.tmpl", targetFileName: "go.mod"},
		{tmplFile: "gitignore.tmpl", targetFileName: ".gitignore"},
		{tmplFile: "main.go.tmpl", targetFileName: "main.go"},
		{tmplFile: "tmpl.go.tmpl", targetFileName: "tmpl.go"},
		{tmplFile: "rawtypes.go.tmpl", targetFileName: "rawtypes.go"},
		{tmplFile: "README.md.tmpl", targetFileName: "README.md"},
	}

	for _, target := range targets {
		rendered, err := tmpl.RenderAsBytes(target.tmplFile, tmplData)
		if err != nil {
			return err
		}

		if err := fs.Add(*jennies.NewFile(target.targetFileName, rendered)); err != nil {
			return err
		}
	}

	if err := fs.Add(*jennies.NewFile("templates/keep.tmpl", []byte("keep"))); err != nil {
		return err
	}

	if err := fs.Write(ctx, opts.outputDir); err != nil {
		return err
	}

	fmt.Printf("Plugin generated in %s\n", opts.outputDir)

	return nil
}

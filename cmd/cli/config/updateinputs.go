package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
	"github.com/grafana/cog/internal/codegen"
	"github.com/grafana/cog/internal/github"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/logs"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type updateInputsOpts struct {
	Write bool
}

func updateInputs(configOpts *options) *cobra.Command {
	opts := &updateInputsOpts{}
	verbosity := 0
	var logger *slog.Logger

	cmd := &cobra.Command{
		Use:   "update-inputs",
		Short: "Inspects inputs and identifies possible schema updates.",
		Long: `Inspects inputs and identifies possible schema updates.

	# List possible updates:

	cog config update-inputs --config config.yaml

	# List possible updates and update configuration file(s):

	cog config update-inputs --config config.yaml --write
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logLevel := new(slog.LevelVar)
			logLevel.Set(slog.LevelWarn)
			// Multiplying the number of occurrences of the `-v` flag by 4 (gap between log levels in slog)
			// allows us to increase the logger's verbosity.
			logLevel.Set(logLevel.Level() - slog.Level(min(verbosity, 3)*4))

			logHandler := logs.NewHandler(os.Stderr, &logs.Options{
				Level: logLevel,
			})
			logger = slog.New(logHandler)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return doUpdateInputs(cmd.Context(), logger, configOpts, opts)
		},
	}

	cmd.Flags().CountVarP(&verbosity, "verbose", "v", "Verbose mode. Multiple -v options increase the verbosity (maximum: 3).")
	cmd.Flags().BoolVar(&opts.Write, "write", opts.Write, "Write the updates back into the configuration file(s).")

	return cmd
}

type updateAction struct {
	Path        string
	SchemaRef   string
	FromSchema  string
	FromVersion string
	ToSchema    string
	ToVersion   string
}

func doUpdateInputs(ctx context.Context, logger *slog.Logger, configOpts *options, opts *updateInputsOpts) error {
	pipelineOpts := []codegen.PipelineOption{
		codegen.Parameters(configOpts.ExtraParameters),
		codegen.Logger(logger),
	}

	pipeline, err := codegen.PipelineFromFile(configOpts.ConfigPath, pipelineOpts...)
	if err != nil {
		return err
	}

	if err := pipeline.LoadUnits(); err != nil {
		return err
	}

	releasesCache := make(map[string]string)
	updates := make([]updateAction, 0, len(pipeline.Inputs))

	for _, input := range pipeline.Inputs {
		l := logger.With(slog.Any("input", stripWorkingDir(input.Source.String())))

		l.Info("checking input for updates")

		if err := inputSupportsUpdates(input); err != nil {
			l.Info("input doesn't support updates", logs.Err(err))
			continue
		}

		schemaURL, urlField, fieldUpdater := inputSchemaUrl(input)

		repoDescriptor := github.RawURLToRepoDescriptor(schemaURL)
		if repoDescriptor == nil {
			l.Warn("could not parse input URL as GitHub raw URL", slog.String("url", schemaURL))
			continue
		}

		repo := repoDescriptor.Owner + "/" + repoDescriptor.Name
		l = l.With(slog.String("repo", repo))

		if !semver.IsValid(repoDescriptor.Ref) {
			l.Warn("invalid semver version", slog.String("version", repoDescriptor.Ref))
			continue
		}

		latestReleaseTag := releasesCache[repo]
		if latestReleaseTag == "" {
			releaseTag, err := github.LatestReleaseTag(ctx, *repoDescriptor)
			if err != nil {
				l.Warn("could not fetch latest release tag", logs.Err(err))
				continue
			}

			releasesCache[repo] = releaseTag
			latestReleaseTag = releaseTag
		}

		if semver.Compare(repoDescriptor.Ref, latestReleaseTag) >= 0 {
			l.Info("no newer release found", slog.String("current_tag", repoDescriptor.Ref), slog.String("latest_tag", latestReleaseTag))
			continue
		}

		newSchemaURL := strings.Replace(schemaURL, repoDescriptor.Ref, latestReleaseTag, 1)
		schemaRef := input.Source.Ref + "." + urlField

		currentSchemas, err := input.LoadSchemas(ctx, logger)
		if err != nil {
			l.Error("could not load current schemas", logs.Err(err))
			continue
		}

		fieldUpdater(newSchemaURL)
		newSchemas, err := input.LoadSchemas(ctx, logger)
		if err != nil {
			l.Error("could not load new schemas", logs.Err(err))
			continue
		}

		identical, err := schemasAreIdentical(currentSchemas, newSchemas)
		if err != nil {
			l.Error("could not check if current and new schemas are identical", logs.Err(err))
			continue
		}

		if identical {
			l.Info("current and new schemas are identical")
			continue
		}

		if opts.Write {
			if err := updateConfigFile(input.Source.Path, schemaRef, newSchemaURL); err != nil {
				l.Error("could write updated config file", logs.Err(err))
				continue
			}
		}

		l.Info("found suitable update", slog.String("from", repoDescriptor.Ref), slog.String("to", latestReleaseTag))

		updates = append(updates, updateAction{
			Path:        input.Source.Path,
			SchemaRef:   schemaRef,
			FromSchema:  schemaURL,
			FromVersion: repoDescriptor.Ref,
			ToSchema:    newSchemaURL,
			ToVersion:   latestReleaseTag,
		})
	}

	renderUpdatesAsMarkdown(updates)

	return nil
}

func renderUpdatesAsMarkdown(updates []updateAction) {
	if len(updates) == 0 {
		return
	}

	fmt.Printf("| File | Ref | Schema | From | To |\n")
	fmt.Printf("|------|-----|--------|------|----|\n")

	for _, entry := range updates {
		config := stripWorkingDir(entry.Path)
		fmt.Printf("| %s | %s | %s | `%s` | `%s` |\n", config, entry.SchemaRef, entry.FromSchema, entry.FromVersion, entry.ToVersion)
	}
}

func stripWorkingDir(input string) string {
	workingDir, _ := os.Getwd()
	stripped := strings.ReplaceAll(input, workingDir, "")
	stripped = strings.TrimPrefix(stripped, string(filepath.Separator))

	return stripped
}

func inputSupportsUpdates(input *codegen.Input) error {
	// no source, no actionable output: we just skip.
	if input.Source == nil {
		return fmt.Errorf("no source")
	}

	url, _, _ := inputSchemaUrl(input)
	if url == "" {
		return fmt.Errorf("no schema URL")
	}

	return nil
}

// returns URL, name of the URL field in the YAML input
func inputSchemaUrl(input *codegen.Input) (string, string, func(newUrl string)) {
	if input.JSONSchema != nil {
		return input.JSONSchema.URL, "jsonschema.url", func(newUrl string) {
			input.JSONSchema.URL = newUrl
		}
	}

	if input.OpenAPI != nil {
		return input.OpenAPI.URL, "openapi.url", func(newUrl string) {
			input.OpenAPI.URL = newUrl
		}
	}

	if input.Cue != nil {
		return input.Cue.URL, "cue.url", func(newUrl string) {
			input.Cue.URL = newUrl
		}
	}

	return "", "", func(newUrl string) {}
}

func updateConfigFile(file string, ref string, value string) error {
	path, err := yaml.PathString(ref)
	if err != nil {
		return fmt.Errorf("could not parse YAML path: %w", err)
	}

	newValue, err := yaml.NewEncoder(nil).EncodeToNode(value)
	if err != nil {
		return fmt.Errorf("could not encode new version as YAML node: %w", err)
	}

	fileAst, err := parser.ParseFile(file, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("could not parse input file: %w", err)
	}

	err = path.ReplaceWithNode(fileAst, newValue)
	if err != nil {
		return fmt.Errorf("could not encode replace YAML node: %w", err)
	}

	err = os.WriteFile(file, []byte(fileAst.String()), 0600)
	if err != nil {
		return fmt.Errorf("could not write updated YAML file: %w", err)
	}

	return nil
}

func schemasAreIdentical(left ir.Schemas, right ir.Schemas) (bool, error) {
	leftJson, err := json.Marshal(left)
	if err != nil {
		return false, fmt.Errorf("could not convert 'left' schemas to JSON: %w", err)
	}

	rightJson, err := json.Marshal(right)
	if err != nil {
		return false, fmt.Errorf("could not convert 'right' schemas to JSON: %w", err)
	}

	return string(leftJson) == string(rightJson), nil
}

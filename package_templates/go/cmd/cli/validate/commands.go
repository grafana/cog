package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-foundation-sdk/go/cmd/cli/tools"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/urfave/cli/v2"
)

type unmarshaller func(decoder *json.Decoder, input []byte, hint string) (any, error)

var unmarshallers = map[string]unmarshaller{
	tools.KindDashboard: unmarshalDashboard,
	tools.KindPanel:     unmarshalPanel,
	tools.KindQuery:     unmarshalQuery,
}

func Command() *cli.Command {
	return &cli.Command{
		Name:        "validate",
		Usage:       "Validate a Grafana resource",
		Description: "Validates a resource JSON from INPUT (or standard input).",
		ArgsUsage:   "[INPUT]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "kind",
				Aliases:  []string{"k"},
				Required: true,
			},
			&cli.StringFlag{
				Name: "hint",
			},
			&cli.BoolFlag{
				Name:    "strict",
				Aliases: []string{"s"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			input, err := tools.ReadFileOrStdin(cCtx.Args().First())
			if err != nil {
				return err
			}

			decoder := json.NewDecoder(bytes.NewBuffer(input))

			if cCtx.Bool("strict") {
				decoder.DisallowUnknownFields()
			}

			_, err = unmarshalKind(decoder, input, cCtx.String("kind"), cCtx.String("hint"))

			return err
		},
	}
}

func unmarshalKind(decoder *json.Decoder, input []byte, kind string, hint string) (any, error) {
	unmarshallerFunc, found := unmarshallers[kind]
	if !found {
		return nil, fmt.Errorf("unknown kind '%s'. Valid kinds are: %s", kind, strings.Join(tools.KnownKinds(), ", "))
	}

	return unmarshallerFunc(decoder, input, hint)
}

func unmarshalDashboard(decoder *json.Decoder, input []byte, _ string) (any, error) {
	dash := dashboard.Dashboard{}
	if err := decoder.Decode(&dash); err != nil {
		return dash, err
	}

	return dash, nil
}

func unmarshalPanel(decoder *json.Decoder, input []byte, _ string) (any, error) {
	panel := dashboard.Panel{}
	if err := decoder.Decode(&panel); err != nil {
		return panel, err
	}

	return panel, nil
}

func unmarshalQuery(decoder *json.Decoder, input []byte, hint string) (any, error) {
	query, err := cog.UnmarshalDataquery(input, hint)
	if err != nil {
		return nil, err
	}

	return query, nil
}

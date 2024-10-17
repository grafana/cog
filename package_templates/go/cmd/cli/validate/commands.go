package validate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-foundation-sdk/go/cmd/cli/tools"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/urfave/cli/v2"
)

type unmarshaller func(input []byte, hint string) (any, error)

var unmarshallers = map[string]unmarshaller{
	tools.KindDashboard: unmarshalDashboard,
	tools.KindPanel:     unmarshalPanel,
	tools.KindQuery:     unmarshalQuery,
}

var strictUnmarshallers = map[string]unmarshaller{
	tools.KindDashboard: unmarshalDashboardStrict,
	tools.KindPanel:     unmarshalPanelStrict,
	tools.KindQuery:     unmarshalQueryStrict,
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
				Usage:    fmt.Sprintf("Supported kinds: %s", strings.Join(tools.KnownKinds(), ", ")),
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

			unmarshallersMap := unmarshallers
			if cCtx.Bool("strict") {
				unmarshallersMap = strictUnmarshallers
			}

			_, err = unmarshalKind(unmarshallersMap, input, cCtx.String("kind"), cCtx.String("hint"))

			return err
		},
	}
}

func unmarshalKind(unmarshallersMap map[string]unmarshaller, input []byte, kind string, hint string) (any, error) {
	unmarshallerFunc, found := unmarshallersMap[kind]
	if !found {
		return nil, fmt.Errorf("unknown kind '%s'. Valid kinds are: %s", kind, strings.Join(tools.KnownKinds(), ", "))
	}

	return unmarshallerFunc(input, hint)
}

func unmarshalDashboard(input []byte, _ string) (any, error) {
	dash := dashboard.Dashboard{}
	if err := json.Unmarshal(input, &dash); err != nil {
		return nil, err
	}

	return dash, nil
}

func unmarshalPanel(input []byte, _ string) (any, error) {
	panel := dashboard.Panel{}
	if err := json.Unmarshal(input, &panel); err != nil {
		return nil, err
	}

	return panel, nil
}

func unmarshalQuery(input []byte, hint string) (any, error) {
	query, err := cog.UnmarshalDataquery(input, hint)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func unmarshalDashboardStrict(input []byte, _ string) (any, error) {
	dash := dashboard.Dashboard{}
	if err := dash.UnmarshalJSONStrict(input); err != nil {
		return nil, err
	}

	return dash, nil
}

func unmarshalPanelStrict(input []byte, _ string) (any, error) {
	panel := dashboard.Panel{}
	if err := panel.UnmarshalJSONStrict(input); err != nil {
		return nil, err
	}

	return panel, nil
}

func unmarshalQueryStrict(input []byte, hint string) (any, error) {
	query, err := cog.StrictUnmarshalDataquery(input, hint)
	if err != nil {
		return nil, err
	}

	return query, nil
}

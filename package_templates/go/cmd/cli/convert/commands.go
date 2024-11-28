package convert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-foundation-sdk/go/cmd/cli/tools"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/urfave/cli/v2"
)

type converter func(input []byte, hint string) (string, error)

func Command() *cli.Command {
	converters := map[string]converter{
		tools.KindDashboard: convertDashboard,
		tools.KindPanel:     convertPanel,
		tools.KindQuery:     convertQuery,
	}

	return &cli.Command{
		Name:        "convert",
		Usage:       "Convert Grafana resources to Go",
		Description: "Converts a Grafana resource JSON from INPUT (or standard input) to Go.",
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
		},
		Action: func(cCtx *cli.Context) error {
			kind := cCtx.String("kind")
			converterFunc, found := converters[kind]
			if !found {
				return fmt.Errorf("unknown kind '%s'. Valid kinds are: %s", kind, strings.Join(tools.KnownKinds(), ", "))
			}

			input, err := tools.ReadFileOrStdin(cCtx.Args().First())
			if err != nil {
				return err
			}

			result, err := converterFunc(input, cCtx.String("hint"))
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}
}

func convertDashboard(input []byte, _ string) (string, error) {
	dash := dashboard.Dashboard{}
	if err := json.Unmarshal(input, &dash); err != nil {
		return "", err
	}

	return dashboard.DashboardConverter(dash), nil
}

func convertPanel(input []byte, _ string) (string, error) {
	panel := dashboard.Panel{}
	if err := json.Unmarshal(input, &panel); err != nil {
		return "", err
	}

	return cog.ConvertPanelToCode(panel, panel.Type), nil
}

func convertQuery(input []byte, hint string) (string, error) {
	query, err := cog.UnmarshalDataquery(input, hint)
	if err != nil {
		return "", err
	}

	return cog.ConvertDataqueryToCode(query), nil
}

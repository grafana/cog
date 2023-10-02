package grafanaplugin

import (
	"github.com/grafana/kindsys"
	"github.com/grafana/grafana/packages/grafana-schema/src/common"
)

kindsys.Composable & {
	name:            "TimeSeries" + "PanelCfg"
	schemaInterface: "PanelCfg"
	lineage: {
		schemas: [{
			version: [0, 0]
			schema: {
				Options: common.OptionsWithTimezones & {
						legend:  common.VizLegendOptions
						tooltip: common.VizTooltipOptions
				}            @cuetsy(kind="interface")
				FieldConfig: common.GraphFieldConfig @cuetsy(kind="interface")
			}
		}]
		name: "TimeSeries" + "PanelCfg"
		lenses: []
	}
}

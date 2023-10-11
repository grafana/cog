package grafanaplugin

import (
	"github.com/grafana/kindsys"
	"github.com/grafana/grafana/packages/grafana-schema/src/common"
)

kindsys.Composable & {
	maturity:        "experimental"
	name:            "Logs" + "PanelCfg"
	schemaInterface: "PanelCfg"
	lineage: {
		schemas: [{
			version: [0, 0]
			schema: {
				Options: {
					showLabels:         bool
					showCommonLabels:   bool
					showTime:           bool
					wrapLogMessage:     bool
					prettifyLogMessage: bool
					enableLogDetails:   bool
					sortOrder:          common.LogsSortOrder
					dedupStrategy:      common.LogsDedupStrategy
				} @cuetsy(kind="interface")
			}
		}]
		name: "Logs" + "PanelCfg"
		lenses: []
	}
}

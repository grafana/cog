package timeseries

import (
	"github.com/grafana/grafana/packages/grafana-schema/src/common"
)

Options: common.OptionsWithTimezones & {
	legend:  common.VizLegendOptions
	tooltip: common.VizTooltipOptions
}

FieldConfig: common.GraphFieldConfig

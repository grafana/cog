package variant_panelcfg_full

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type OptionsDataSourceModel struct {
     TimeseriesOption types.String `tfsdk:"timeseries_option"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type FieldConfigDataSourceModel struct {
     TimeseriesFieldConfigOption types.String `tfsdk:"timeseries_field_config_option"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


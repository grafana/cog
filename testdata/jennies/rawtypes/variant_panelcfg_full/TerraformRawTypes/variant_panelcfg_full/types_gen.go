package variant_panelcfg_full

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type Options struct {
 TimeseriesOption types.String `tfsdk:"timeseries_option"`
 }

type FieldConfig struct {
 TimeseriesFieldConfigOption types.String `tfsdk:"timeseries_field_config_option"`
 }


package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type VariantPanelcfgFullOptionsDataSourceModel struct {
  
     TimeseriesOption types.String `tfsdk:"timeseries_option"`
  
}

type VariantPanelcfgFullFieldConfigDataSourceModel struct {
  
     TimeseriesFieldConfigOption types.String `tfsdk:"timeseries_field_config_option"`
  
}


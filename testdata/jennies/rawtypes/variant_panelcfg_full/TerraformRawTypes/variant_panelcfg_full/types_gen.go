package variant_panelcfg_full

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type OptionsModel struct {
	TimeseriesOption types.String `tfsdk:"timeseries_option"`
}
var OptionsType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"timeseries_option": types.StringType,
	},
}
var OptionsSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"timeseries_option": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type FieldConfigModel struct {
	TimeseriesFieldConfigOption types.String `tfsdk:"timeseries_field_config_option"`
}
var FieldConfigType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"timeseries_field_config_option": types.StringType,
	},
}
var FieldConfigSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"timeseries_field_config_option": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


package variant_panelcfg_full

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Options struct {
 TimeseriesOption types.String `tfsdk:"timeseries_option"`
 }

type FieldConfig struct {
 TimeseriesFieldConfigOption types.String `tfsdk:"timeseries_field_config_option"`
 }

var OptionsAttributes = map[string]schema.Attribute{
"timeseries_option": schema.StringAttribute{
 Required: true,
},

}

var FieldConfigAttributes = map[string]schema.Attribute{
"timeseries_field_config_option": schema.StringAttribute{
 Required: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
"options": schema.SingleNestedAttribute{
Required: true,
Attributes: OptionsAttributes,
},
"field_config": schema.SingleNestedAttribute{
Required: true,
Attributes: FieldConfigAttributes,
},
}
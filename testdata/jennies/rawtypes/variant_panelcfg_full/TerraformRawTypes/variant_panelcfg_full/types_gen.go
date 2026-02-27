package variant_panelcfg_full

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type Options struct {
 TimeseriesOption types.String `tfsdk:"timeseries_option"`
 }

type FieldConfig struct {
 TimeseriesFieldConfigOption types.String `tfsdk:"timeseries_field_config_option"`
 }

var SpecAttributes = map[string]schema.Attribute{
"options": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"timeseries_option": types.StringType,
},
},
"field_config": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"timeseries_field_config_option": types.StringType,
},
},
}
package variant_panelcfg_only_options

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type OptionsModel struct {
	Content types.String `tfsdk:"content"`
}
var OptionsType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"content": types.StringType,
	},
}
var OptionsSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"content": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


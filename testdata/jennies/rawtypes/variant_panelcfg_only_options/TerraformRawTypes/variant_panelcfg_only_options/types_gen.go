package variant_panelcfg_only_options

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type Options struct {
 Content types.String `tfsdk:"content"`
 }

var SpecAttributes = map[string]schema.Attribute{
"options": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"content": types.StringType,
},
},
}
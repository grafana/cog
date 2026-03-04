package variant_panelcfg_only_options

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Options struct {
 Content types.String `tfsdk:"content"`
 }

var SpecAttributes = map[string]schema.Attribute{
"options": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"content": schema.StringAttribute{
 Required: true,
},

},
},
}
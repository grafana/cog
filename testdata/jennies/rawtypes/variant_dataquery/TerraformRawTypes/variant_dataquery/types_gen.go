package variant_dataquery

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type Query struct {
 Expr types.String `tfsdk:"expr"`
Instant types.Bool `tfsdk:"instant"`
 }

var SpecAttributes = map[string]schema.Attribute{
"query": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"expr": types.StringType,
"instant": types.BoolType,
},
},
}
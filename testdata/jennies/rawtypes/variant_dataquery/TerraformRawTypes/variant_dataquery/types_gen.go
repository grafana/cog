package variant_dataquery

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Query struct {
 Expr types.String `tfsdk:"expr"`
Instant types.Bool `tfsdk:"instant"`
 }

var QueryAttributes = map[string]schema.Attribute{
"expr": schema.StringAttribute{
 Required: true,
},

"instant": schema.BoolAttribute{
 Optional: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
"query": schema.SingleNestedAttribute{
Required: true,
Attributes: QueryAttributes,
},
}
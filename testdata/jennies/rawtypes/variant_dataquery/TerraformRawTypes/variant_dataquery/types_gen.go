package variant_dataquery

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type QueryModel struct {
	Expr types.String `tfsdk:"expr"`
	Instant types.Bool `tfsdk:"instant"`
}
var QueryType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"expr": types.StringType,
		"instant": types.BoolType,
	},
}
var QuerySchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"expr": schema.StringAttribute{
	Required: true,
},
	"instant": schema.BoolAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}


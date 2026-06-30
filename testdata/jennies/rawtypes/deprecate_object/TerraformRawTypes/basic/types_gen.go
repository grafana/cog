package basic

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type SomeStructModel struct {
	FieldString types.String `tfsdk:"field_string"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_string": types.StringType,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	DeprecationMessage: "This object is deprecated, use NewStruct instead.",
	Attributes: map[string]schema.Attribute{
	"field_string": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


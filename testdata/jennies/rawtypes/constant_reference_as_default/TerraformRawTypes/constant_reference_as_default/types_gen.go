package constant_reference_as_default

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const ConstantRefStringModel = "AString"
var ConstantRefStringType = types.StringType


type MyStructModel struct {
}
var MyStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
	},
}
var MyStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
},
}


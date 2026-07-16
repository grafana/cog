package reference_of_reference

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type MyStructModel struct {
	Field types.Object `tfsdk:"field"`
}
var MyStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field": OtherStructType,
	},
}
var MyStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
	"field": schema.SingleNestedBlock{
		Attributes: AnotherStructSchema.Attributes,
		Blocks: AnotherStructSchema.Blocks,
		Validators: AnotherStructSchema.Validators,
	},
},
}

type OtherStructModel = types.Object
var OtherStructType = AnotherStructType


type AnotherStructModel struct {
	A types.String `tfsdk:"a"`
}
var AnotherStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"a": types.StringType,
	},
}
var AnotherStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"a": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


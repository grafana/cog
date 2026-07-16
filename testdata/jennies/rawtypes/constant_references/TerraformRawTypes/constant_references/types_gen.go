package constant_references

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

type EnumModel types.String
var EnumType = types.StringType


type ParentStructModel struct {
	MyEnum types.String `tfsdk:"my_enum"`
}
var ParentStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"my_enum": types.StringType,
	},
}
var ParentStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"my_enum": schema.StringAttribute{
	Required: true,
	Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},
},
},
	Blocks: map[string]schema.Block{
},
}

type StructModel struct {
	MyValue types.String `tfsdk:"my_value"`
	MyEnum types.String `tfsdk:"my_enum"`
}
var StructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"my_value": types.StringType,
		"my_enum": types.StringType,
	},
}
var StructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"my_value": schema.StringAttribute{
	Required: true,
},
	"my_enum": schema.StringAttribute{
	Required: true,
	Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},
},
},
	Blocks: map[string]schema.Block{
},
}

type StructAModel struct {
}
var StructAType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
	},
}
var StructASchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
},
}

type StructBModel struct {
	MyValue types.String `tfsdk:"my_value"`
}
var StructBType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"my_value": types.StringType,
	},
}
var StructBSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"my_value": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


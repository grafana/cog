package struct_optional_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

type SomeStructModel struct {
	FieldRef types.Object `tfsdk:"field_ref"`
	FieldString types.String `tfsdk:"field_string"`
	Operator types.String `tfsdk:"operator"`
	FieldArrayOfStrings types.List `tfsdk:"field_array_of_strings"`
	FieldAnonymousStruct types.Object `tfsdk:"field_anonymous_struct"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_ref": SomeOtherStructType,
		"field_string": types.StringType,
		"operator": types.StringType,
		"field_array_of_strings": types.ListType{
	ElemType: types.StringType,
},
		"field_anonymous_struct": StructOptionalFieldsSomeStructFieldAnonymousStructType,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_string": schema.StringAttribute{
	Optional: true,
},
	"operator": schema.StringAttribute{
	Required: true,
	Validators: []validator.String{
stringvalidator.OneOf(">", "<"),
},
},
	"field_array_of_strings": schema.ListAttribute{
	ElementType: types.StringType,
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
	"field_ref": schema.SingleNestedBlock{
		Attributes: SomeOtherStructSchema.Attributes,
		Blocks: SomeOtherStructSchema.Blocks,
		Validators: SomeOtherStructSchema.Validators,
	},
	"field_anonymous_struct": schema.SingleNestedBlock{
		Attributes: StructOptionalFieldsSomeStructFieldAnonymousStructSchema.Attributes,
		Blocks: StructOptionalFieldsSomeStructFieldAnonymousStructSchema.Blocks,
		Validators: StructOptionalFieldsSomeStructFieldAnonymousStructSchema.Validators,
	},
},
}

type SomeOtherStructModel struct {
	FieldAny types.String `tfsdk:"field_any"`
}
var SomeOtherStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_any": types.StringType,
	},
}
var SomeOtherStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_any": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type StructOptionalFieldsSomeStructFieldAnonymousStructModel struct {
	FieldAny types.String `tfsdk:"field_any"`
}
var StructOptionalFieldsSomeStructFieldAnonymousStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_any": types.StringType,
	},
}
var StructOptionalFieldsSomeStructFieldAnonymousStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_any": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type SomeStructOperatorModel types.String
var SomeStructOperatorType = types.StringType



package struct_complex_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// This struct does things.
type SomeStructModel struct {
	FieldRef types.Object `tfsdk:"field_ref"`
	FieldDisjunctionOfScalars types.Object `tfsdk:"field_disjunction_of_scalars"`
	FieldMixedDisjunction types.Object `tfsdk:"field_mixed_disjunction"`
	FieldDisjunctionWithNull types.String `tfsdk:"field_disjunction_with_null"`
	Operator types.String `tfsdk:"operator"`
	FieldArrayOfStrings types.List `tfsdk:"field_array_of_strings"`
	FieldMapOfStringToString types.Map `tfsdk:"field_map_of_string_to_string"`
	FieldAnonymousStruct types.Object `tfsdk:"field_anonymous_struct"`
	FieldRefToConstant types.String `tfsdk:"field_ref_to_constant"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_ref": SomeOtherStructType,
		"field_disjunction_of_scalars": StringOrBoolType,
		"field_mixed_disjunction": StringOrSomeOtherStructType,
		"field_disjunction_with_null": types.StringType,
		"operator": types.StringType,
		"field_array_of_strings": types.ListType{
	ElemType: types.StringType,
},
		"field_map_of_string_to_string": types.MapType{
	ElemType: types.StringType,
},
		"field_anonymous_struct": StructComplexFieldsSomeStructFieldAnonymousStructType,
		"field_ref_to_constant": ConnectionPathType,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Description: "This struct does things.",
	MarkdownDescription: "This struct does things.",
	Attributes: map[string]schema.Attribute{
	"field_disjunction_with_null": schema.StringAttribute{
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
	Required: true,
},
	"field_map_of_string_to_string": schema.MapAttribute{
	ElementType: types.StringType,
	Required: true,
},
	"field_ref_to_constant": schema.StringAttribute{
	Required: true,
	Default: stringdefault.StaticString("straight"),
	Computed: true,
},
},
	Blocks: map[string]schema.Block{
	"field_ref": schema.SingleNestedBlock{
		Attributes: SomeOtherStructSchema.Attributes,
		Blocks: SomeOtherStructSchema.Blocks,
		Validators: SomeOtherStructSchema.Validators,
	},
	"field_disjunction_of_scalars": schema.SingleNestedBlock{
		Attributes: StringOrBoolSchema.Attributes,
		Blocks: StringOrBoolSchema.Blocks,
		Validators: StringOrBoolSchema.Validators,
	},
	"field_mixed_disjunction": schema.SingleNestedBlock{
		Attributes: StringOrSomeOtherStructSchema.Attributes,
		Blocks: StringOrSomeOtherStructSchema.Blocks,
		Validators: StringOrSomeOtherStructSchema.Validators,
	},
	"field_anonymous_struct": schema.SingleNestedBlock{
		Attributes: StructComplexFieldsSomeStructFieldAnonymousStructSchema.Attributes,
		Blocks: StructComplexFieldsSomeStructFieldAnonymousStructSchema.Blocks,
		Validators: StructComplexFieldsSomeStructFieldAnonymousStructSchema.Validators,
	},
},
}

const ConnectionPathModel = "straight"
var ConnectionPathType = types.StringType


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

type StructComplexFieldsSomeStructFieldAnonymousStructModel struct {
	FieldAny types.String `tfsdk:"field_any"`
}
var StructComplexFieldsSomeStructFieldAnonymousStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_any": types.StringType,
	},
}
var StructComplexFieldsSomeStructFieldAnonymousStructSchema = schema.SingleNestedBlock{
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


type StringOrBoolModel struct {
	String types.String `tfsdk:"string"`
	Bool types.Bool `tfsdk:"bool"`
}
var StringOrBoolType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"string": types.StringType,
		"bool": types.BoolType,
	},
}
var StringOrBoolSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"string": schema.StringAttribute{
	Optional: true,
},
	"bool": schema.BoolAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
	Validators: []validator.Object{
		(1, "string", "bool"),
	},
}

type StringOrSomeOtherStructModel struct {
	String types.String `tfsdk:"string"`
	SomeOtherStruct types.Object `tfsdk:"some_other_struct"`
}
var StringOrSomeOtherStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"string": types.StringType,
		"some_other_struct": SomeOtherStructType,
	},
}
var StringOrSomeOtherStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"string": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
	"some_other_struct": schema.SingleNestedBlock{
		Attributes: SomeOtherStructSchema.Attributes,
		Blocks: SomeOtherStructSchema.Blocks,
		Validators: SomeOtherStructSchema.Validators,
	},
},
}


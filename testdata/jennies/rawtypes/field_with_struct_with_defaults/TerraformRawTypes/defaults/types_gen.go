package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type NestedStructModel struct {
	StringVal types.String `tfsdk:"string_val"`
	IntVal types.Int64 `tfsdk:"int_val"`
}
var NestedStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"string_val": types.StringType,
		"int_val": types.Int64Type,
	},
}
var NestedStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"string_val": schema.StringAttribute{
	Required: true,
},
	"int_val": schema.Int64Attribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type StructModel struct {
	AllFields types.Object `tfsdk:"all_fields"`
	PartialFields types.Object `tfsdk:"partial_fields"`
	EmptyFields types.Object `tfsdk:"empty_fields"`
	ComplexField types.Object `tfsdk:"complex_field"`
	PartialComplexField types.Object `tfsdk:"partial_complex_field"`
}
var StructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"all_fields": NestedStructType,
		"partial_fields": NestedStructType,
		"empty_fields": NestedStructType,
		"complex_field": DefaultsStructComplexFieldType,
		"partial_complex_field": DefaultsStructPartialComplexFieldType,
	},
}
var StructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
	"all_fields": schema.SingleNestedBlock{
		Attributes: NestedStructSchema.Attributes,
		Blocks: NestedStructSchema.Blocks,
		Validators: NestedStructSchema.Validators,
	},
	"partial_fields": schema.SingleNestedBlock{
		Attributes: NestedStructSchema.Attributes,
		Blocks: NestedStructSchema.Blocks,
		Validators: NestedStructSchema.Validators,
	},
	"empty_fields": schema.SingleNestedBlock{
		Attributes: NestedStructSchema.Attributes,
		Blocks: NestedStructSchema.Blocks,
		Validators: NestedStructSchema.Validators,
	},
	"complex_field": schema.SingleNestedBlock{
		Attributes: DefaultsStructComplexFieldSchema.Attributes,
		Blocks: DefaultsStructComplexFieldSchema.Blocks,
		Validators: DefaultsStructComplexFieldSchema.Validators,
	},
	"partial_complex_field": schema.SingleNestedBlock{
		Attributes: DefaultsStructPartialComplexFieldSchema.Attributes,
		Blocks: DefaultsStructPartialComplexFieldSchema.Blocks,
		Validators: DefaultsStructPartialComplexFieldSchema.Validators,
	},
},
}

type DefaultsStructComplexFieldNestedModel struct {
	NestedVal types.String `tfsdk:"nested_val"`
}
var DefaultsStructComplexFieldNestedType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"nested_val": types.StringType,
	},
}
var DefaultsStructComplexFieldNestedSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"nested_val": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type DefaultsStructComplexFieldModel struct {
	Uid types.String `tfsdk:"uid"`
	Nested types.Object `tfsdk:"nested"`
	Array types.List `tfsdk:"array"`
}
var DefaultsStructComplexFieldType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uid": types.StringType,
		"nested": DefaultsStructComplexFieldNestedType,
		"array": types.ListType{
	ElemType: types.StringType,
},
	},
}
var DefaultsStructComplexFieldSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"uid": schema.StringAttribute{
	Required: true,
},
	"array": schema.ListAttribute{
	ElementType: types.StringType,
	Required: true,
},
},
	Blocks: map[string]schema.Block{
	"nested": schema.SingleNestedBlock{
		Attributes: DefaultsStructComplexFieldNestedSchema.Attributes,
		Blocks: DefaultsStructComplexFieldNestedSchema.Blocks,
		Validators: DefaultsStructComplexFieldNestedSchema.Validators,
	},
},
}

type DefaultsStructPartialComplexFieldModel struct {
	Uid types.String `tfsdk:"uid"`
	IntVal types.Int64 `tfsdk:"int_val"`
}
var DefaultsStructPartialComplexFieldType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"uid": types.StringType,
		"int_val": types.Int64Type,
	},
}
var DefaultsStructPartialComplexFieldSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"uid": schema.StringAttribute{
	Required: true,
},
	"int_val": schema.Int64Attribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type NestedStruct struct {
 StringVal types.String `tfsdk:"stringVal"`
IntVal types.Int64 `tfsdk:"intVal"`
 }

type Struct struct {
 AllFields NestedStruct `tfsdk:"allFields"`
PartialFields NestedStruct `tfsdk:"partialFields"`
EmptyFields NestedStruct `tfsdk:"emptyFields"`
ComplexField DefaultsStructComplexField `tfsdk:"complexField"`
PartialComplexField DefaultsStructPartialComplexField `tfsdk:"partialComplexField"`
 }

type DefaultsStructComplexFieldNested struct {
 NestedVal types.String `tfsdk:"nestedVal"`
 }

type DefaultsStructComplexField struct {
 Uid types.String `tfsdk:"uid"`
Nested DefaultsStructComplexFieldNested `tfsdk:"nested"`
Array types.List `tfsdk:"array"`
 }

type DefaultsStructPartialComplexField struct {
 Uid types.String `tfsdk:"uid"`
IntVal types.Int64 `tfsdk:"intVal"`
 }

var NestedStructAttributes = map[string]schema.Attribute{
"string_val": schema.StringAttribute{
 Required: true,
},

"int_val": schema.Int64Attribute{
 Required: true,
},

}

var StructAttributes = map[string]schema.Attribute{
"all_fields": schema.SingleNestedAttribute{
Required: true,
Attributes: NestedStructAttributes,
},

"partial_fields": schema.SingleNestedAttribute{
Required: true,
Attributes: NestedStructAttributes,
},

"empty_fields": schema.SingleNestedAttribute{
Required: true,
Attributes: NestedStructAttributes,
},

"complex_field": schema.SingleNestedAttribute{
Required: true,
Attributes: DefaultsStructComplexFieldAttributes,
},

"partial_complex_field": schema.SingleNestedAttribute{
Required: true,
Attributes: DefaultsStructPartialComplexFieldAttributes,
},

}

var DefaultsStructComplexFieldNestedAttributes = map[string]schema.Attribute{
"nested_val": schema.StringAttribute{
 Required: true,
},

}

var DefaultsStructComplexFieldAttributes = map[string]schema.Attribute{
"uid": schema.StringAttribute{
 Required: true,
},

"nested": schema.SingleNestedAttribute{
Required: true,
Attributes: DefaultsStructComplexFieldNestedAttributes,
},

"array": schema.ListAttribute{
 ElementType: types.StringType,
},

}

var DefaultsStructPartialComplexFieldAttributes = map[string]schema.Attribute{
"uid": schema.StringAttribute{
 Required: true,
},

"int_val": schema.Int64Attribute{
 Required: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
"nested_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: NestedStructAttributes,
},
"struct": schema.SingleNestedAttribute{
Required: true,
Attributes: StructAttributes,
},
"defaults_struct_complex_field_nested": schema.SingleNestedAttribute{
Required: true,
Attributes: DefaultsStructComplexFieldNestedAttributes,
},
"defaults_struct_complex_field": schema.SingleNestedAttribute{
Required: true,
Attributes: DefaultsStructComplexFieldAttributes,
},
"defaults_struct_partial_complex_field": schema.SingleNestedAttribute{
Required: true,
Attributes: DefaultsStructPartialComplexFieldAttributes,
},
}
package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
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

var SpecAttributes = map[string]schema.Attribute{
"nestedstruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"stringVal": types.StringType,
"intVal": types.Int64Type,
},
},
"struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"allFields": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"stringVal": types.StringType,
"intVal": types.Int64Type,
},
},
"partialFields": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"stringVal": types.StringType,
"intVal": types.Int64Type,
},
},
"emptyFields": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"stringVal": types.StringType,
"intVal": types.Int64Type,
},
},
"complexField": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"uid": types.StringType,
"nested": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"nestedVal": types.StringType,
},
},
"array": types.ListType{
 ElemType: types.StringType,
},
},
},
"partialComplexField": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"uid": types.StringType,
"intVal": types.Int64Type,
},
},
},
},
"defaultsstructcomplexfieldnested": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"nestedVal": types.StringType,
},
},
"defaultsstructcomplexfield": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"uid": types.StringType,
"nested": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"nestedVal": types.StringType,
},
},
"array": types.ListType{
 ElemType: types.StringType,
},
},
},
"defaultsstructpartialcomplexfield": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"uid": types.StringType,
"intVal": types.Int64Type,
},
},
}
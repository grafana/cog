package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


package basic

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

// This
// is
// a
// comment
type SomeStruct struct {
 // Anything can go in there.
// Really, anything.
FieldAny types.Object `tfsdk:"FieldAny"`
FieldBool types.Bool `tfsdk:"FieldBool"`
FieldBytes types.String `tfsdk:"FieldBytes"`
FieldString types.String `tfsdk:"FieldString"`
FieldStringWithConstantValue types.String `tfsdk:"FieldStringWithConstantValue"`
FieldFloat32 types.Float32 `tfsdk:"FieldFloat32"`
FieldFloat64 types.Float64 `tfsdk:"FieldFloat64"`
FieldUint8 types.Number `tfsdk:"FieldUint8"`
FieldUint16 types.Number `tfsdk:"FieldUint16"`
FieldUint32 types.Int32 `tfsdk:"FieldUint32"`
FieldUint64 types.Int64 `tfsdk:"FieldUint64"`
FieldInt8 types.Number `tfsdk:"FieldInt8"`
FieldInt16 types.Number `tfsdk:"FieldInt16"`
FieldInt32 types.Int32 `tfsdk:"FieldInt32"`
FieldInt64 types.Int64 `tfsdk:"FieldInt64"`
 }

var SpecAttributes = map[string]schema.Attribute{
"somestruct": types.ObjectAttributes{
Required: true,
Description: `
This
is
a
comment
`,
,AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
"FieldBool": types.BoolType,
"FieldBytes": types.StringType,
"FieldString": types.StringType,
"FieldStringWithConstantValue": types.StringType,
"FieldFloat32": types.Float32Type,
"FieldFloat64": types.Float64Type,
"FieldUint8": types.NumberType,
"FieldUint16": types.NumberType,
"FieldUint32": types.Int32Type,
"FieldUint64": types.Int64Type,
"FieldInt8": types.NumberType,
"FieldInt16": types.NumberType,
"FieldInt32": types.Int32Type,
"FieldInt64": types.Int64Type,
},
}
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
"someStruct": schema.ObjectAttribute{
Required: true,
Description: `
This
is
a
comment
`,
AttributeTypes: map[string]attr.Type{
"field_any": types.ObjectType{},
"field_bool": types.BoolType,
"field_bytes": types.StringType,
"field_string": types.StringType,
"field_string_with_constant_value": types.StringType,
"field_float32": types.Float32Type,
"field_float64": types.Float64Type,
"field_uint8": types.NumberType,
"field_uint16": types.NumberType,
"field_uint32": types.Int32Type,
"field_uint64": types.Int64Type,
"field_int8": types.NumberType,
"field_int16": types.NumberType,
"field_int32": types.Int32Type,
"field_int64": types.Int64Type,
},
},
}
package basic

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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

var SomeStructAttributes = map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

"field_bool": schema.BoolAttribute{
 Required: true,
},

"field_bytes": schema.StringAttribute{
 Required: true,
},

"field_string": schema.StringAttribute{
 Required: true,
},

"field_string_with_constant_value": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("auto"),
},

"field_float32": schema.Float32Attribute{
 Required: true,
},

"field_float64": schema.Float64Attribute{
 Required: true,
},

"field_uint8": schema.NumberAttribute{
 Required: true,
},

"field_uint16": schema.NumberAttribute{
 Required: true,
},

"field_uint32": schema.Int32Attribute{
 Required: true,
},

"field_uint64": schema.Int64Attribute{
 Required: true,
},

"field_int8": schema.NumberAttribute{
 Required: true,
},

"field_int16": schema.NumberAttribute{
 Required: true,
},

"field_int32": schema.Int32Attribute{
 Required: true,
},

"field_int64": schema.Int64Attribute{
 Required: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Description: `
This
is
a
comment
`,
Attributes: SomeStructAttributes,
},
}
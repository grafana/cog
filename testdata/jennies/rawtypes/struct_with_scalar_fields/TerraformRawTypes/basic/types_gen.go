package basic

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


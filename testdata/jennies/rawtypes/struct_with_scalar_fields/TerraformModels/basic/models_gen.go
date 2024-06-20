package basic

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldAny types.Object `tfsdk:"field_any"`
     FieldBool types.Bool `tfsdk:"field_bool"`
     FieldBytes types.String `tfsdk:"field_bytes"`
     FieldString types.String `tfsdk:"field_string"`
     FieldStringWithConstantValue types.String `tfsdk:"field_string_with_constant_value"`
     FieldFloat32 types.Float64 `tfsdk:"field_float32"`
     FieldFloat64 types.Float64 `tfsdk:"field_float64"`
     FieldUint8 types.Int64 `tfsdk:"field_uint8"`
     FieldUint16 types.Int64 `tfsdk:"field_uint16"`
     FieldUint32 types.Int64 `tfsdk:"field_uint32"`
     FieldUint64 types.Int64 `tfsdk:"field_uint64"`
     FieldInt8 types.Int64 `tfsdk:"field_int8"`
     FieldInt16 types.Int64 `tfsdk:"field_int16"`
     FieldInt32 types.Int64 `tfsdk:"field_int32"`
     FieldInt64 types.Int64 `tfsdk:"field_int64"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


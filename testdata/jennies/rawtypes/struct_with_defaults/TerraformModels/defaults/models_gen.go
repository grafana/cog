package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldBool types.Bool `tfsdk:"field_bool"`
     FieldString types.String `tfsdk:"field_string"`
     FieldStringWithConstantValue types.String `tfsdk:"field_string_with_constant_value"`
     FieldFloat32 types.Float64 `tfsdk:"field_float32"`
     FieldInt32 types.Int64 `tfsdk:"field_int32"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type SomeStruct struct {
 FieldBool types.Bool `tfsdk:"fieldBool"`
FieldString types.String `tfsdk:"fieldString"`
FieldStringWithConstantValue types.String `tfsdk:"FieldStringWithConstantValue"`
FieldFloat32 types.Float32 `tfsdk:"FieldFloat32"`
FieldInt32 types.Int32 `tfsdk:"FieldInt32"`
 }


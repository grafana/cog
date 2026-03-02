package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type SomeStruct struct {
 FieldBool types.Bool `tfsdk:"fieldBool"`
FieldString types.String `tfsdk:"fieldString"`
FieldStringWithConstantValue types.String `tfsdk:"FieldStringWithConstantValue"`
FieldFloat32 types.Float32 `tfsdk:"FieldFloat32"`
FieldInt32 types.Int32 `tfsdk:"FieldInt32"`
 }

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"field_bool": types.BoolType,
"field_string": types.StringType,
"field_string_with_constant_value": types.StringType,
"field_float32": types.Float32Type,
"field_int32": types.Int32Type,
},
},
}
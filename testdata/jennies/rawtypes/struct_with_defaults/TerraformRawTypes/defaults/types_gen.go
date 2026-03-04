package defaults

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	booldefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	float32default "github.com/hashicorp/terraform-plugin-framework/resource/schema/float32default"
	int32default "github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
)

type SomeStruct struct {
 FieldBool types.Bool `tfsdk:"fieldBool"`
FieldString types.String `tfsdk:"fieldString"`
FieldStringWithConstantValue types.String `tfsdk:"FieldStringWithConstantValue"`
FieldFloat32 types.Float32 `tfsdk:"FieldFloat32"`
FieldInt32 types.Int32 `tfsdk:"FieldInt32"`
 }

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_bool": schema.BoolAttribute{
 Required: true,
Default: booldefault.StaticBool(true)
},

"field_string": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("foo")
},

"field_string_with_constant_value": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("auto"),
},

"field_float32": schema.Float32Attribute{
 Required: true,
Default: float32default.StaticFloat32(42.42)
},

"field_int32": schema.Int32Attribute{
 Required: true,
Default: int32default.StaticInt32(42)
},

},
},
}
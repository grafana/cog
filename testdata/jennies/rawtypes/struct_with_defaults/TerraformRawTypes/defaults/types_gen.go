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
"somestruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldBool": types.BoolType,
"fieldString": types.StringType,
"fieldStringWithConstantValue": types.StringType,
"fieldFloat32": types.Float32Type,
"fieldInt32": types.Int32Type,
},
},
}
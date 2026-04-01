package constant_reference_as_default

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

const ConstantRefString = "AString"

type MyStruct struct {
 AString types.String `tfsdk:"aString"`
OptString types.String `tfsdk:"optString"`
 }

var SpecAttributes = map[string]schema.Attribute{
"constant_ref_string": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("AString"),
},
"my_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"a_string": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("AString"),
},

"opt_string": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("AString"),
},

},
},
}
package constant_reference_as_default

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

const ConstantRefString = "AString"

type MyStruct struct {
 AString types.String `tfsdk:"aString"`
OptString types.String `tfsdk:"optString"`
 }

var SpecAttributes = map[string]schema.Attribute{
"constantRefString": schema.StringAttribute{
 Required: true,

},
"myStruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"a_string": types.StringType,
"opt_string": types.StringType,
},
},
}
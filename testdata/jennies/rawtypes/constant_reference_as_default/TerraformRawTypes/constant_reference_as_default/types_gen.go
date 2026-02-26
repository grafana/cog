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
"constantrefstring": schema.StringAttribute{
 Required: true
 
}"mystruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"aString": unknown,
"optString": unknown,
},
}
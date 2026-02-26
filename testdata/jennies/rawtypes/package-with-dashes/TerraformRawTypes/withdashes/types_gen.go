package withdashes

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

// Refresh rate or disabled.
type RefreshRate = StringOrBool

type StringOrBool struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
 }

var SpecAttributes = map[string]schema.Attribute{
"somestruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
}
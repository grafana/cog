package withdashes

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

var SomeStructAttributes = map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

}

var StringOrBoolAttributes = map[string]schema.Attribute{
"string": schema.StringAttribute{
 Optional: true,
},

"bool": schema.BoolAttribute{
 Optional: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: SomeStructAttributes,
},
}
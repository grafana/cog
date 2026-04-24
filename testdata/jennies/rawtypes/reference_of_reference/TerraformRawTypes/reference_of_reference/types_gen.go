package reference_of_reference

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type MyStruct struct {
 Field OtherStruct `tfsdk:"field"`
 }

type OtherStruct = AnotherStruct

type AnotherStruct struct {
 A types.String `tfsdk:"a"`
 }

var SpecAttributes = map[string]schema.Attribute{
"my_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"a": schema.StringAttribute{
 Required: true,
},

},
},

},
},
"another_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"a": schema.StringAttribute{
 Required: true,
},

},
},
}
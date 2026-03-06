package refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type RefToSomeStruct = SomeStruct

type RefToSomeStructFromOtherPackage = unknown

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
}
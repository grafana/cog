package refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type RefToSomeStruct = SomeStruct

type RefToSomeStructFromOtherPackage = unknown

var SpecAttributes = map[string]schema.Attribute{
"somestruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
}
package arrays

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// List of tags, maybe?
type ArrayOfStrings types.List

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type ArrayOfRefs []SomeStruct

type ArrayOfArrayOfNumbers types.List

var SpecAttributes = map[string]schema.Attribute{
"array_of_strings": schema.ListAttribute{
 ElementType: types.StringType,
},
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
"array_of_refs": schema.ListNestedAttribute{
NestedObject: schema.NestedAttributeObject {
Attributes: map[string]schema.Attribute {
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
},
"array_of_array_of_numbers": schema.ListAttribute{
 ElementType: types.ListType{
 ElemType: types.Int64Type,
},
},
}
package arrays

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

// List of tags, maybe?
type ArrayOfStrings types.List

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type ArrayOfRefs []SomeStruct

type ArrayOfArrayOfNumbers types.List

var SpecAttributes = map[string]schema.Attribute{
"arrayofstrings": schema.ListAttribute{
 ElementType: types.StringType,
},
"somestruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
"arrayofrefs": schema.ListAttribute{
 ElementType: types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
},
"arrayofarrayofnumbers": schema.ListAttribute{
 ElementType: types.ListType{
 ElemType: types.Int64Type,
},
},
}
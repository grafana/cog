package maps

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

// String to... something.
type MapOfStringToAny types.Map

type MapOfStringToString types.Map

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type MapOfStringToRef types.Map

type MapOfStringToMapOfStringToBool types.Map

var SpecAttributes = map[string]schema.Attribute{
"map_of_string_to_any": schema.ListMapAttribute{
 ElementType: types.ObjectType{},
},
"map_of_string_to_string": schema.ListMapAttribute{
 ElementType: types.StringType,
},
"some_struct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"field_any": types.ObjectType{},
},
},
"map_of_string_to_ref": schema.ListMapAttribute{
 ElementType: types.ObjectType{
 AttrTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
},
"map_of_string_to_map_of_string_to_bool": schema.ListMapAttribute{
 ElementType: types.MapType{
 ElemType: types.BoolType,
},
},
}
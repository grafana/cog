package maps

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
"map_of_string_to_any": schema.MapAttribute{
 ElementType: types.DynamicType,
},
"map_of_string_to_string": schema.MapAttribute{
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
"map_of_string_to_ref": schema.MapNestedAttribute{
NestedObject: schema.NestedAttributeObject {
Attributes: map[string]schema.Attribute {
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
},
"map_of_string_to_map_of_string_to_bool": schema.MapAttribute{
 ElementType: types.MapType{
 ElemType: types.BoolType,
},
},
}
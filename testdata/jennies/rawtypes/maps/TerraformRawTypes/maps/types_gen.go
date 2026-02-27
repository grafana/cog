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
"mapOfStringToAny": schema.ListMapAttribute{
 ElementType: types.ObjectType{},
},
"mapOfStringToString": schema.ListMapAttribute{
 ElementType: types.StringType,
},
"someStruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
"mapOfStringToRef": schema.ListMapAttribute{
 ElementType: types.ObjectType{
 AttrTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
},
"mapOfStringToMapOfStringToBool": schema.ListMapAttribute{
 ElementType: types.MapType{
 ElemType: types.BoolType,
},
},
}
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
"mapofstringtoany": types.MapAttribute{
 
}"mapofstringtostring": types.MapAttribute{
 
}"somestruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
"mapofstringtoref": types.MapAttribute{
 
}"mapofstringtomapofstringtobool": types.MapAttribute{
 
}}
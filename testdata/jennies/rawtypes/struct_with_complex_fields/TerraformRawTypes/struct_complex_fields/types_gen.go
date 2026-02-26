package struct_complex_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

// This struct does things.
type SomeStruct struct {
 FieldRef SomeOtherStruct `tfsdk:"FieldRef"`
FieldDisjunctionOfScalars StringOrBool `tfsdk:"FieldDisjunctionOfScalars"`
FieldMixedDisjunction StringOrSomeOtherStruct `tfsdk:"FieldMixedDisjunction"`
FieldDisjunctionWithNull types.String `tfsdk:"FieldDisjunctionWithNull"`
Operator types.String `tfsdk:"Operator"`
FieldArrayOfStrings types.List `tfsdk:"FieldArrayOfStrings"`
FieldMapOfStringToString types.Map `tfsdk:"FieldMapOfStringToString"`
FieldAnonymousStruct StructComplexFieldsSomeStructFieldAnonymousStruct `tfsdk:"FieldAnonymousStruct"`
FieldRefToConstant ConnectionPath `tfsdk:"fieldRefToConstant"`
 }

const ConnectionPath = "straight"

type SomeOtherStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type StructComplexFieldsSomeStructFieldAnonymousStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }



type StringOrBool struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
 }

type StringOrSomeOtherStruct struct {
 String types.String `tfsdk:"String"`
SomeOtherStruct SomeOtherStruct `tfsdk:"SomeOtherStruct"`
 }

var SpecAttributes = map[string]schema.Attribute{
"somestruct": types.ObjectAttributes{
Required: true,
Description: `
This struct does things.
`,
,AttributeTypes: map[string]attr.Type{
"FieldRef": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
,
"FieldDisjunctionOfScalars": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"String": types.StringType,
"Bool": types.BoolType,
},
,
"FieldMixedDisjunction": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"String": types.StringType,
"SomeOtherStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
,
},
,
"FieldDisjunctionWithNull": types.StringType,
"Operator": unknown,
"FieldArrayOfStrings": types.ListType{
 ElemType: types.StringType,
},
"FieldMapOfStringToString": types.MapType{
 ElemType: types.StringType,
},
"FieldAnonymousStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
,
"fieldRefToConstant": types.StringType,
},
"connectionpath": schema.StringAttribute{
 Required: true
 
}"someotherstruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
"structcomplexfieldssomestructfieldanonymousstruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
"somestructoperator": "stringorbool": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"String": types.StringType,
"Bool": types.BoolType,
},
"stringorsomeotherstruct": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"String": types.StringType,
"SomeOtherStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
,
},
}
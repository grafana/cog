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
"somestruct": schema.ObjectAttribute{
Required: true,
Description: `
This struct does things.
`,
AttributeTypes: map[string]attr.Type{
"fieldRef": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
"fieldDisjunctionOfScalars": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"String": types.StringType,
"Bool": types.BoolType,
},
},
"fieldMixedDisjunction": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"String": types.StringType,
"SomeOtherStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
},
},
"fieldDisjunctionWithNull": types.StringType,
"operator": types.StringType,
"fieldArrayOfStrings": types.ListType{
 ElemType: types.StringType,
},
"fieldMapOfStringToString": types.MapType{
 ElemType: types.StringType,
},
"fieldAnonymousStruct": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"FieldAny": types.ObjectType{},
},
},
"fieldRefToConstant": types.StringType,
},
},
"connectionpath": schema.StringAttribute{
 Required: true,

},
"someotherstruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
"structcomplexfieldssomestructfieldanonymousstruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldAny": types.ObjectType{},
},
},
}
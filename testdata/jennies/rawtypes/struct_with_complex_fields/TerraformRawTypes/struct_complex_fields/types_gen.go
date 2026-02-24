package struct_complex_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


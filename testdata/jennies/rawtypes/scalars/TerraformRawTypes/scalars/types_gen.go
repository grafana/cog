package scalars

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

const ConstTypeString = "foo"

type ScalarTypeAny types.Object

type ScalarTypeBool types.Bool

type ScalarTypeBytes types.String

type ScalarTypeString types.String

type ScalarTypeFloat32 types.Float32

type ScalarTypeFloat64 types.Float64

type ScalarTypeUint8 types.Number

type ScalarTypeUint16 types.Number

type ScalarTypeUint32 types.Int32

type ScalarTypeUint64 types.Int64

type ScalarTypeInt8 types.Number

type ScalarTypeInt16 types.Number

type ScalarTypeInt32 types.Int32

type ScalarTypeInt64 types.Int64

var SpecAttributes = map[string]schema.Attribute{
"consttypestring": schema.StringAttribute{
 Required: true, 
},
"scalartypeany": schema.ObjectAttribute{
 Required: true, 
},
"scalartypebool": schema.BoolAttribute{
 Required: true, 
},
"scalartypebytes": schema.StringAttribute{
 Required: true, 
},
"scalartypestring": schema.StringAttribute{
 Required: true, 
},
"scalartypefloat32": schema.Float32Attribute{
 Required: true, 
},
"scalartypefloat64": schema.Float64Attribute{
 Required: true, 
},
"scalartypeuint8": schema.NumberAttribute{
 Required: true, 
},
"scalartypeuint16": schema.NumberAttribute{
 Required: true, 
},
"scalartypeuint32": schema.Int32Attribute{
 Required: true, 
},
"scalartypeuint64": schema.Int64Attribute{
 Required: true, 
},
"scalartypeint8": schema.NumberAttribute{
 Required: true, 
},
"scalartypeint16": schema.NumberAttribute{
 Required: true, 
},
"scalartypeint32": schema.Int32Attribute{
 Required: true, 
},
"scalartypeint64": schema.Int64Attribute{
 Required: true, 
},
}
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
"constTypeString": schema.StringAttribute{
 Required: true,

},
"scalarTypeAny": schema.ObjectAttribute{
 Required: true, 
},
"scalarTypeBool": schema.BoolAttribute{
 Required: true, 
},
"scalarTypeBytes": schema.StringAttribute{
 Required: true,

},
"scalarTypeString": schema.StringAttribute{
 Required: true,

},
"scalarTypeFloat32": schema.Float32Attribute{
 Required: true, 
},
"scalarTypeFloat64": schema.Float64Attribute{
 Required: true, 
},
"scalarTypeUint8": schema.NumberAttribute{
 Required: true, 
},
"scalarTypeUint16": schema.NumberAttribute{
 Required: true, 
},
"scalarTypeUint32": schema.Int32Attribute{
 Required: true, 
},
"scalarTypeUint64": schema.Int64Attribute{
 Required: true, 
},
"scalarTypeInt8": schema.NumberAttribute{
 Required: true, 
},
"scalarTypeInt16": schema.NumberAttribute{
 Required: true, 
},
"scalarTypeInt32": schema.Int32Attribute{
 Required: true, 
},
"scalarTypeInt64": schema.Int64Attribute{
 Required: true, 
},
}
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
"const_type_string": schema.StringAttribute{
 Required: true,

},
"scalar_type_any": schema.ObjectAttribute{
 Required: true, 
},
"scalar_type_bool": schema.BoolAttribute{
 Required: true, 
},
"scalar_type_bytes": schema.StringAttribute{
 Required: true,

},
"scalar_type_string": schema.StringAttribute{
 Required: true,

},
"scalar_type_float32": schema.Float32Attribute{
 Required: true, 
},
"scalar_type_float64": schema.Float64Attribute{
 Required: true, 
},
"scalar_type_uint8": schema.NumberAttribute{
 Required: true, 
},
"scalar_type_uint16": schema.NumberAttribute{
 Required: true, 
},
"scalar_type_uint32": schema.Int32Attribute{
 Required: true, 
},
"scalar_type_uint64": schema.Int64Attribute{
 Required: true, 
},
"scalar_type_int8": schema.NumberAttribute{
 Required: true, 
},
"scalar_type_int16": schema.NumberAttribute{
 Required: true, 
},
"scalar_type_int32": schema.Int32Attribute{
 Required: true, 
},
"scalar_type_int64": schema.Int64Attribute{
 Required: true, 
},
}
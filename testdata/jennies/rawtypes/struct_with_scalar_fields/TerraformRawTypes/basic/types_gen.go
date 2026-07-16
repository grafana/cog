package basic

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// This
// is
// a
// comment
type SomeStructModel struct {
	// Anything can go in there.
	// Really, anything.
	FieldAny types.String `tfsdk:"field_any"`
	FieldBool types.Bool `tfsdk:"field_bool"`
	FieldBytes types.String `tfsdk:"field_bytes"`
	FieldString types.String `tfsdk:"field_string"`
	FieldStringWithConstantValue types.String `tfsdk:"field_string_with_constant_value"`
	FieldFloat32 types.Float32 `tfsdk:"field_float32"`
	FieldFloat64 types.Float64 `tfsdk:"field_float64"`
	FieldUint8 types.Number `tfsdk:"field_uint8"`
	FieldUint16 types.Number `tfsdk:"field_uint16"`
	FieldUint32 types.Int32 `tfsdk:"field_uint32"`
	FieldUint64 types.Int64 `tfsdk:"field_uint64"`
	FieldInt8 types.Number `tfsdk:"field_int8"`
	FieldInt16 types.Number `tfsdk:"field_int16"`
	FieldInt32 types.Int32 `tfsdk:"field_int32"`
	FieldInt64 types.Int64 `tfsdk:"field_int64"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_any": types.StringType,
		"field_bool": types.BoolType,
		"field_bytes": types.StringType,
		"field_string": types.StringType,
		"field_string_with_constant_value": types.StringType,
		"field_float32": types.Float32Type,
		"field_float64": types.Float64Type,
		"field_uint8": types.NumberType,
		"field_uint16": types.NumberType,
		"field_uint32": types.Int32Type,
		"field_uint64": types.Int64Type,
		"field_int8": types.NumberType,
		"field_int16": types.NumberType,
		"field_int32": types.Int32Type,
		"field_int64": types.Int64Type,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Description: "Thisisacomment",
	MarkdownDescription: "Thisisacomment",
	Attributes: map[string]schema.Attribute{
	"field_any": schema.StringAttribute{
	Required: true,
	Description: "Anything can go in there.Really, anything.",
	MarkdownDescription: "Anything can go in there.Really, anything.",
},
	"field_bool": schema.BoolAttribute{
	Required: true,
},
	"field_bytes": schema.StringAttribute{
	Required: true,
},
	"field_string": schema.StringAttribute{
	Required: true,
},
	"field_string_with_constant_value": schema.StringAttribute{
	Required: true,
	Default: stringdefault.StaticString("auto"),
	Computed: true,
},
	"field_float32": schema.Float32Attribute{
	Required: true,
},
	"field_float64": schema.Float64Attribute{
	Required: true,
},
	"field_uint8": schema.NumberAttribute{
	Required: true,
},
	"field_uint16": schema.NumberAttribute{
	Required: true,
},
	"field_uint32": schema.Int32Attribute{
	Required: true,
},
	"field_uint64": schema.Int64Attribute{
	Required: true,
},
	"field_int8": schema.NumberAttribute{
	Required: true,
},
	"field_int16": schema.NumberAttribute{
	Required: true,
},
	"field_int32": schema.Int32Attribute{
	Required: true,
},
	"field_int64": schema.Int64Attribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


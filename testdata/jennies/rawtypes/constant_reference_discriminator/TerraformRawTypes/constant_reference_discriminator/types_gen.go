package constant_reference_discriminator

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type LayoutWithValue = GridLayoutUsingValueOrRowsLayoutUsingValue

type GridLayoutUsingValue struct {
 Kind types.String `tfsdk:"kind"`
GridLayoutProperty types.String `tfsdk:"gridLayoutProperty"`
 }

type RowsLayoutUsingValue struct {
 Kind types.String `tfsdk:"kind"`
RowsLayoutProperty types.String `tfsdk:"rowsLayoutProperty"`
 }

type LayoutWithoutValue = GridLayoutWithoutValueOrRowsLayoutWithoutValue

type GridLayoutWithoutValue struct {
 Kind types.String `tfsdk:"kind"`
GridLayoutProperty types.String `tfsdk:"gridLayoutProperty"`
 }

type RowsLayoutWithoutValue struct {
 Kind types.String `tfsdk:"kind"`
RowsLayoutProperty types.String `tfsdk:"rowsLayoutProperty"`
 }

const GridLayoutKindType = "GridLayout"

const RowsLayoutKindType = "RowsLayout"

type GridLayoutUsingValueOrRowsLayoutUsingValue struct {
 GridLayoutUsingValue GridLayoutUsingValue `tfsdk:"GridLayoutUsingValue"`
RowsLayoutUsingValue RowsLayoutUsingValue `tfsdk:"RowsLayoutUsingValue"`
 }

type GridLayoutWithoutValueOrRowsLayoutWithoutValue struct {
 GridLayoutWithoutValue GridLayoutWithoutValue `tfsdk:"GridLayoutWithoutValue"`
RowsLayoutWithoutValue RowsLayoutWithoutValue `tfsdk:"RowsLayoutWithoutValue"`
 }

var SpecAttributes = map[string]schema.Attribute{
"gridLayoutUsingValue": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"kind": types.StringType,
"gridLayoutProperty": types.StringType,
},
},
"rowsLayoutUsingValue": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"kind": types.StringType,
"rowsLayoutProperty": types.StringType,
},
},
"gridLayoutWithoutValue": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"kind": types.StringType,
"gridLayoutProperty": types.StringType,
},
},
"rowsLayoutWithoutValue": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"kind": types.StringType,
"rowsLayoutProperty": types.StringType,
},
},
"gridLayoutKindType": schema.StringAttribute{
 Required: true,

},
"rowsLayoutKindType": schema.StringAttribute{
 Required: true,

},
}
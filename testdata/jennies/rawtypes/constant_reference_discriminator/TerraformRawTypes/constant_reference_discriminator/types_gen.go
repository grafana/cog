package constant_reference_discriminator

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
"grid_layout_using_value": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"kind": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("GridLayout"),
},

"grid_layout_property": schema.StringAttribute{
 Required: true,
},

},
},
"rows_layout_using_value": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"kind": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("RowsLayout"),
},

"rows_layout_property": schema.StringAttribute{
 Required: true,
},

},
},
"grid_layout_without_value": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"kind": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("GridLayout"),
},

"grid_layout_property": schema.StringAttribute{
 Required: true,
},

},
},
"rows_layout_without_value": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"kind": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("RowsLayout"),
},

"rows_layout_property": schema.StringAttribute{
 Required: true,
},

},
},
"grid_layout_kind_type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("GridLayout"),
},
"rows_layout_kind_type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("RowsLayout"),
},
}
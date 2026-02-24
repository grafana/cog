package constant_reference_discriminator

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


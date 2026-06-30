package constant_reference_discriminator

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type LayoutWithValueModel = types.Object
var LayoutWithValueType = GridLayoutUsingValueOrRowsLayoutUsingValueType


type GridLayoutUsingValueModel struct {
	GridLayoutProperty types.String `tfsdk:"grid_layout_property"`
}
var GridLayoutUsingValueType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"grid_layout_property": types.StringType,
	},
}
var GridLayoutUsingValueSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"grid_layout_property": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type RowsLayoutUsingValueModel struct {
	RowsLayoutProperty types.String `tfsdk:"rows_layout_property"`
}
var RowsLayoutUsingValueType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"rows_layout_property": types.StringType,
	},
}
var RowsLayoutUsingValueSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"rows_layout_property": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type LayoutWithoutValueModel = types.Object
var LayoutWithoutValueType = GridLayoutWithoutValueOrRowsLayoutWithoutValueType


type GridLayoutWithoutValueModel struct {
	GridLayoutProperty types.String `tfsdk:"grid_layout_property"`
}
var GridLayoutWithoutValueType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"grid_layout_property": types.StringType,
	},
}
var GridLayoutWithoutValueSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"grid_layout_property": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type RowsLayoutWithoutValueModel struct {
	RowsLayoutProperty types.String `tfsdk:"rows_layout_property"`
}
var RowsLayoutWithoutValueType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"rows_layout_property": types.StringType,
	},
}
var RowsLayoutWithoutValueSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"rows_layout_property": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

const GridLayoutKindTypeModel = "GridLayout"
var GridLayoutKindTypeType = types.StringType


const RowsLayoutKindTypeModel = "RowsLayout"
var RowsLayoutKindTypeType = types.StringType


type GridLayoutUsingValueOrRowsLayoutUsingValueModel struct {
	GridLayoutUsingValue types.Object `tfsdk:"grid_layout_using_value"`
	RowsLayoutUsingValue types.Object `tfsdk:"rows_layout_using_value"`
}
var GridLayoutUsingValueOrRowsLayoutUsingValueType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"grid_layout_using_value": GridLayoutUsingValueType,
		"rows_layout_using_value": RowsLayoutUsingValueType,
	},
}
var GridLayoutUsingValueOrRowsLayoutUsingValueSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
	"grid_layout_using_value": schema.SingleNestedBlock{
		Attributes: GridLayoutUsingValueSchema.Attributes,
		Blocks: GridLayoutUsingValueSchema.Blocks,
		Validators: []validator.Object{
			("grid_layout_property"),
		},
	},
	"rows_layout_using_value": schema.SingleNestedBlock{
		Attributes: RowsLayoutUsingValueSchema.Attributes,
		Blocks: RowsLayoutUsingValueSchema.Blocks,
		Validators: []validator.Object{
			("rows_layout_property"),
		},
	},
},
	Validators: []validator.Object{
		(1, "grid_layout_using_value", "rows_layout_using_value"),
	},
}

type GridLayoutWithoutValueOrRowsLayoutWithoutValueModel struct {
	GridLayoutWithoutValue types.Object `tfsdk:"grid_layout_without_value"`
	RowsLayoutWithoutValue types.Object `tfsdk:"rows_layout_without_value"`
}
var GridLayoutWithoutValueOrRowsLayoutWithoutValueType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"grid_layout_without_value": GridLayoutWithoutValueType,
		"rows_layout_without_value": RowsLayoutWithoutValueType,
	},
}
var GridLayoutWithoutValueOrRowsLayoutWithoutValueSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
	"grid_layout_without_value": schema.SingleNestedBlock{
		Attributes: GridLayoutWithoutValueSchema.Attributes,
		Blocks: GridLayoutWithoutValueSchema.Blocks,
		Validators: []validator.Object{
			("grid_layout_property"),
		},
	},
	"rows_layout_without_value": schema.SingleNestedBlock{
		Attributes: RowsLayoutWithoutValueSchema.Attributes,
		Blocks: RowsLayoutWithoutValueSchema.Blocks,
		Validators: []validator.Object{
			("rows_layout_property"),
		},
	},
},
	Validators: []validator.Object{
		(1, "grid_layout_without_value", "rows_layout_without_value"),
	},
}


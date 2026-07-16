package dashboard

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type DashboardModel struct {
	Title types.String `tfsdk:"title"`
	Panels types.List `tfsdk:"panels"`
}
var DashboardType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"title": types.StringType,
		"panels": types.ListType{
	ElemType: PanelType,
},
	},
}
var DashboardSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"title": schema.StringAttribute{
	Required: true,
},
	"panels": schema.ListAttribute{
	ElementType: PanelType,
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type DataSourceRefModel struct {
	Type types.String `tfsdk:"type"`
	Uid types.String `tfsdk:"uid"`
}
var DataSourceRefType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"uid": types.StringType,
	},
}
var DataSourceRefSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"type": schema.StringAttribute{
	Optional: true,
},
	"uid": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type FieldConfigSourceModel struct {
	Defaults types.Object `tfsdk:"defaults"`
}
var FieldConfigSourceType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"defaults": FieldConfigType,
	},
}
var FieldConfigSourceSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
	"defaults": schema.SingleNestedBlock{
		Attributes: FieldConfigSchema.Attributes,
		Blocks: FieldConfigSchema.Blocks,
		Validators: FieldConfigSchema.Validators,
	},
},
}

type FieldConfigModel struct {
	Unit types.String `tfsdk:"unit"`
	Custom types.String `tfsdk:"custom"`
}
var FieldConfigType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"unit": types.StringType,
		"custom": types.StringType,
	},
}
var FieldConfigSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"unit": schema.StringAttribute{
	Optional: true,
},
	"custom": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type PanelModel struct {
	Title types.String `tfsdk:"title"`
	Type types.String `tfsdk:"type"`
	Datasource types.Object `tfsdk:"datasource"`
	Options types.String `tfsdk:"options"`
	Targets types.List `tfsdk:"targets"`
	FieldConfig types.Object `tfsdk:"field_config"`
}
var PanelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"title": types.StringType,
		"type": types.StringType,
		"datasource": DataSourceRefType,
		"options": types.StringType,
		"targets": types.ListType{
	ElemType: unknown,
},
		"field_config": FieldConfigSourceType,
	},
}
var PanelSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"title": schema.StringAttribute{
	Required: true,
},
	"type": schema.StringAttribute{
	Required: true,
},
	"options": schema.StringAttribute{
	Optional: true,
},
	"targets": schema.ListAttribute{
	ElementType: unknown,
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
	"datasource": schema.SingleNestedBlock{
		Attributes: DataSourceRefSchema.Attributes,
		Blocks: DataSourceRefSchema.Blocks,
		Validators: DataSourceRefSchema.Validators,
	},
	"field_config": schema.SingleNestedBlock{
		Attributes: FieldConfigSourceSchema.Attributes,
		Blocks: FieldConfigSourceSchema.Blocks,
		Validators: FieldConfigSourceSchema.Validators,
	},
},
}


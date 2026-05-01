package dashboard

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Dashboard struct {
 Title types.String `tfsdk:"title"`
Panels []Panel `tfsdk:"panels"`
 }

type DataSourceRef struct {
 Type types.String `tfsdk:"type"`
Uid types.String `tfsdk:"uid"`
 }

type FieldConfigSource struct {
 Defaults FieldConfig `tfsdk:"defaults"`
 }

type FieldConfig struct {
 Unit types.String `tfsdk:"unit"`
Custom types.Object `tfsdk:"custom"`
 }

type Panel struct {
 Title types.String `tfsdk:"title"`
Type types.String `tfsdk:"type"`
Datasource DataSourceRef `tfsdk:"datasource"`
Options types.Object `tfsdk:"options"`
Targets types.List `tfsdk:"targets"`
FieldConfig FieldConfigSource `tfsdk:"fieldConfig"`
 }

var DashboardAttributes = map[string]schema.Attribute{
"title": schema.StringAttribute{
 Required: true,
},

"panels": schema.ListNestedAttribute{
NestedObject: schema.NestedAttributeObject{
Attributes: PanelAttributes,
},
},

}

var DataSourceRefAttributes = map[string]schema.Attribute{
"type": schema.StringAttribute{
 Optional: true,
},

"uid": schema.StringAttribute{
 Optional: true,
},

}

var FieldConfigSourceAttributes = map[string]schema.Attribute{
"defaults": schema.SingleNestedAttribute{
Optional: true,
Attributes: FieldConfigAttributes,
},

}

var FieldConfigAttributes = map[string]schema.Attribute{
"unit": schema.StringAttribute{
 Optional: true,
},

"custom": schema.ObjectAttribute{
 Optional: true,
},

}

var PanelAttributes = map[string]schema.Attribute{
"title": schema.StringAttribute{
 Required: true,
},

"type": schema.StringAttribute{
 Required: true,
},

"datasource": schema.SingleNestedAttribute{
Optional: true,
Attributes: DataSourceRefAttributes,
},

"options": schema.ObjectAttribute{
 Optional: true,
},

"targets": schema.ListAttribute{
 ElementType: unknown,
},

"field_config": schema.SingleNestedAttribute{
Optional: true,
Attributes: FieldConfigSourceAttributes,
},

}

var SpecAttributes = map[string]schema.Attribute{
"dashboard": schema.SingleNestedAttribute{
Required: true,
Attributes: DashboardAttributes,
},
"data_source_ref": schema.SingleNestedAttribute{
Required: true,
Attributes: DataSourceRefAttributes,
},
"field_config_source": schema.SingleNestedAttribute{
Required: true,
Attributes: FieldConfigSourceAttributes,
},
"field_config": schema.SingleNestedAttribute{
Required: true,
Attributes: FieldConfigAttributes,
},
"panel": schema.SingleNestedAttribute{
Required: true,
Attributes: PanelAttributes,
},
}
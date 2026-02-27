package dashboard

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
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

var SpecAttributes = map[string]schema.Attribute{
"dashboard": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"title": types.StringType,
"panels": types.ListType{
 ElemType: types.ObjectType{
 AttrTypes: map[string]attr.Type{
"title": types.StringType,
"type": types.StringType,
"datasource": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"type": types.StringType,
"uid": types.StringType,
},
},
"options": types.ObjectType{},
"targets": types.ListType{
 ElemType: unknown,
},
"fieldConfig": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"defaults": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"unit": types.StringType,
"custom": types.ObjectType{},
},
},
},
},
},
},
},
},
},
"data_source_ref": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"type": types.StringType,
"uid": types.StringType,
},
},
"field_config_source": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"defaults": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"unit": types.StringType,
"custom": types.ObjectType{},
},
},
},
},
"field_config": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"unit": types.StringType,
"custom": types.ObjectType{},
},
},
"panel": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"title": types.StringType,
"type": types.StringType,
"datasource": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"type": types.StringType,
"uid": types.StringType,
},
},
"options": types.ObjectType{},
"targets": types.ListType{
 ElemType: unknown,
},
"field_config": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"defaults": types.ObjectType{
 AttrTypes: map[string]attr.Type{
"unit": types.StringType,
"custom": types.ObjectType{},
},
},
},
},
},
},
}
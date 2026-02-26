package dashboard

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type DashboardDataSourceModel struct {
  
     Title types.String `tfsdk:"title"`
  
}

type DashboardDataSourceRefDataSourceModel struct {
  
     Type types.String `tfsdk:"type"`
  
     Uid types.String `tfsdk:"uid"`
  
}

type DashboardFieldConfigSourceDataSourceModel struct {
  
}

type DashboardFieldConfigDataSourceModel struct {
  
     Unit types.String `tfsdk:"unit"`
  
     Custom types.Object `tfsdk:"custom"`
  
}

type DashboardPanelDataSourceModel struct {
  
     Title types.String `tfsdk:"title"`
  
     Type types.String `tfsdk:"type"`
  
     Options types.Object `tfsdk:"options"`
  
}


package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type WithDashesSomeStructDataSourceModel struct {
  
     FieldAny types.Object `tfsdk:"field_any"`
  
}

type WithDashesStringOrBoolDataSourceModel struct {
  
     String types.String `tfsdk:"string"`
  
     Bool types.Bool `tfsdk:"bool"`
  
}


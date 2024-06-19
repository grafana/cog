package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type VariantDataqueryQueryDataSourceModel struct {
  
     Expr types.String `tfsdk:"expr"`
  
     Instant types.Bool `tfsdk:"instant"`
  
}


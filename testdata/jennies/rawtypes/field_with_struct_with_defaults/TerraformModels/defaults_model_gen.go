package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type DefaultsNestedStructDataSourceModel struct {
  
     StringVal types.String `tfsdk:"string_val"`
  
     IntVal types.Int64 `tfsdk:"int_val"`
  
}

type DefaultsStructDataSourceModel struct {
  
}


package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type StructOptionalFieldsSomeStructDataSourceModel struct {
  
     FieldString types.String `tfsdk:"field_string"`
  
}

type StructOptionalFieldsSomeOtherStructDataSourceModel struct {
  
     FieldAny types.Object `tfsdk:"field_any"`
  
}


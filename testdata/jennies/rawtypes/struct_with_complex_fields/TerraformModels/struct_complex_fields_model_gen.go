package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type StructComplexFieldsSomeStructDataSourceModel struct {
  
     FieldDisjunctionWithNull types.String `tfsdk:"field_disjunction_with_null"`
  
}

type StructComplexFieldsSomeOtherStructDataSourceModel struct {
  
     FieldAny types.Object `tfsdk:"field_any"`
  
}

type StructComplexFieldsStringOrBoolDataSourceModel struct {
  
     String types.String `tfsdk:"string"`
  
     Bool types.Bool `tfsdk:"bool"`
  
}

type StructComplexFieldsStringOrSomeOtherStructDataSourceModel struct {
  
     String types.String `tfsdk:"string"`
  
}


package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)



type DisjunctionsSomeStructDataSourceModel struct {
  
     Type types.String `tfsdk:"type"`
  
     FieldAny types.Object `tfsdk:"field_any"`
  
}

type DisjunctionsSomeOtherStructDataSourceModel struct {
  
     Type types.String `tfsdk:"type"`
  
     Foo types.String `tfsdk:"foo"`
  
}

type DisjunctionsYetAnotherStructDataSourceModel struct {
  
     Type types.String `tfsdk:"type"`
  
     Bar types.Int64 `tfsdk:"bar"`
  
}

type DisjunctionsStringOrBoolDataSourceModel struct {
  
     String types.String `tfsdk:"string"`
  
     Bool types.Bool `tfsdk:"bool"`
  
}

type DisjunctionsBoolOrSomeStructDataSourceModel struct {
  
     Bool types.Bool `tfsdk:"bool"`
  
}

type DisjunctionsSomeStructOrSomeOtherStructOrYetAnotherStructDataSourceModel struct {
  
}


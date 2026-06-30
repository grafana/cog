package arrays

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// List of tags, maybe?
type ArrayOfStringsModel = types.List
var ArrayOfStringsType = types.ListType{
	ElemType: types.StringType,
}


type SomeStructModel struct {
	FieldAny types.String `tfsdk:"field_any"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_any": types.StringType,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_any": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type ArrayOfRefsModel = types.List
var ArrayOfRefsType = types.ListType{
	ElemType: SomeStructType,
}


type ArrayOfArrayOfNumbersModel = types.List
var ArrayOfArrayOfNumbersType = types.ListType{
	ElemType: types.ListType{
	ElemType: types.Int64Type,
},
}



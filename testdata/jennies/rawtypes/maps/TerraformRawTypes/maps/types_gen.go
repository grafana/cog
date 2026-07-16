package maps

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// String to... something.
type MapOfStringToAnyModel = types.Map
var MapOfStringToAnyType = types.MapType{
	ElemType: types.StringType,
}


type MapOfStringToStringModel = types.Map
var MapOfStringToStringType = types.MapType{
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

type MapOfStringToRefModel = types.Map
var MapOfStringToRefType = types.MapType{
	ElemType: SomeStructType,
}


type MapOfStringToMapOfStringToBoolModel = types.Map
var MapOfStringToMapOfStringToBoolType = types.MapType{
	ElemType: types.MapType{
	ElemType: types.BoolType,
},
}



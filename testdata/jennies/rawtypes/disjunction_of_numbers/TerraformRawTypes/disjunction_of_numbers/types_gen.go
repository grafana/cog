package disjunction_of_numbers

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type NumbersModel = types.Object
var NumbersType = Int64OrFloat64OrFloat32Type


type Int64OrFloat64OrFloat32Model struct {
	Int64 types.Int64 `tfsdk:"int64"`
	Float64 types.Float64 `tfsdk:"float64"`
	Float32 types.Float32 `tfsdk:"float32"`
}
var Int64OrFloat64OrFloat32Type = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"int64": types.Int64Type,
		"float64": types.Float64Type,
		"float32": types.Float32Type,
	},
}
var Int64OrFloat64OrFloat32Schema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"int64": schema.Int64Attribute{
	Optional: true,
},
	"float64": schema.Float64Attribute{
	Optional: true,
},
	"float32": schema.Float32Attribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
	Validators: []validator.Object{
		(1, "int64", "float64", "float32"),
	},
}


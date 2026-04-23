package disjunction_of_numbers

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Numbers = Int64OrFloat64OrFloat32

type Int64OrFloat64OrFloat32 struct {
 Int64 types.Int64 `tfsdk:"Int64"`
Float64 types.Float64 `tfsdk:"Float64"`
Float32 types.Float32 `tfsdk:"Float32"`
 }

var Int64OrFloat64OrFloat32Attributes = map[string]schema.Attribute{
"int64": schema.Int64Attribute{
 Optional: true,
},

"float64": schema.Float64Attribute{
 Optional: true,
},

"float32": schema.Float32Attribute{
 Optional: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
}
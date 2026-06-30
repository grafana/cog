package disjunctions_of_refs_without_discriminator

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type DisjunctionWithoutDiscriminatorModel = types.String
var DisjunctionWithoutDiscriminatorType = types.StringType


type TypeAModel struct {
	FieldA types.String `tfsdk:"field_a"`
}
var TypeAType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_a": types.StringType,
	},
}
var TypeASchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_a": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type TypeBModel struct {
	FieldB types.Int64 `tfsdk:"field_b"`
}
var TypeBType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_b": types.Int64Type,
	},
}
var TypeBSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_b": schema.Int64Attribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


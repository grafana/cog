package withdashes

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

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

// Refresh rate or disabled.
type RefreshRateModel = types.Object
var RefreshRateType = StringOrBoolType


type StringOrBoolModel struct {
	String types.String `tfsdk:"string"`
	Bool types.Bool `tfsdk:"bool"`
}
var StringOrBoolType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"string": types.StringType,
		"bool": types.BoolType,
	},
}
var StringOrBoolSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"string": schema.StringAttribute{
	Optional: true,
},
	"bool": schema.BoolAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
	Validators: []validator.Object{
		(1, "string", "bool"),
	},
}


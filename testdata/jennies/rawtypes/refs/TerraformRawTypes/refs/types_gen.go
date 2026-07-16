package refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

type RefToSomeStructModel = types.Object
var RefToSomeStructType = SomeStructType


type RefToSomeStructFromOtherPackageModel = could not resolve ref 'otherpkg.SomeDistantStruct'
var RefToSomeStructFromOtherPackageType = unknown



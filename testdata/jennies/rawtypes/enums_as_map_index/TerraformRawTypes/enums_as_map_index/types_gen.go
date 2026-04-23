package enums_as_map_index

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)





type SomeStruct struct {
 Data types.Map `tfsdk:"data"`
 }

type SomeStructWithDefaultEnum struct {
 Data types.Map `tfsdk:"data"`
 }

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"data": schema.MapAttribute{
 ElementType: types.StringType,
},

},
},
"some_struct_with_default_enum": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"data": schema.MapAttribute{
 ElementType: types.StringType,
},

},
},
}
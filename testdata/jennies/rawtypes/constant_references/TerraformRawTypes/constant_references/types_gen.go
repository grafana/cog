package constant_references

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)



type ParentStruct struct {
 MyEnum types.String `tfsdk:"myEnum"`
 }

type Struct struct {
 MyValue types.String `tfsdk:"myValue"`
MyEnum types.String `tfsdk:"myEnum"`
 }

type StructA struct {
 MyEnum types.String `tfsdk:"myEnum"`
Other types.String `tfsdk:"other"`
 }

type StructB struct {
 MyEnum types.String `tfsdk:"myEnum"`
MyValue types.String `tfsdk:"myValue"`
 }

var SpecAttributes = map[string]schema.Attribute{
"parent_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"my_enum": schema.StringAttribute{
 Required: true,
},

},
},
"struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"my_value": schema.StringAttribute{
 Required: true,
},

"my_enum": schema.StringAttribute{
 Required: true,
},

},
},
"struct_a": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"my_enum": schema.StringAttribute{
 Required: true,
},

"other": schema.StringAttribute{
 Required: true,
},

},
},
"struct_b": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"my_enum": schema.StringAttribute{
 Required: true,
},

"my_value": schema.StringAttribute{
 Required: true,
},

},
},
}
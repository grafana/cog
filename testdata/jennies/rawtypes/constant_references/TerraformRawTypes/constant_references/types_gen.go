package constant_references

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

var ParentStructAttributes = map[string]schema.Attribute{
"my_enum": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},

},

}

var StructAttributes = map[string]schema.Attribute{
"my_value": schema.StringAttribute{
 Required: true,
},

"my_enum": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},

},

}

var StructAAttributes = map[string]schema.Attribute{
"my_enum": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},

},

"other": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},

},

}

var StructBAttributes = map[string]schema.Attribute{
"my_enum": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("ValueA", "ValueB", "ValueC"),
},

},

"my_value": schema.StringAttribute{
 Required: true,
},

}

var SpecAttributes = map[string]schema.Attribute{
"parent_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: ParentStructAttributes,
},
"struct": schema.SingleNestedAttribute{
Required: true,
Attributes: StructAttributes,
},
"struct_a": schema.SingleNestedAttribute{
Required: true,
Attributes: StructAAttributes,
},
"struct_b": schema.SingleNestedAttribute{
Required: true,
Attributes: StructBAttributes,
},
}
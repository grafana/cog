package struct_optional_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type SomeStruct struct {
 FieldRef SomeOtherStruct `tfsdk:"FieldRef"`
FieldString types.String `tfsdk:"FieldString"`
Operator types.String `tfsdk:"Operator"`
FieldArrayOfStrings types.List `tfsdk:"FieldArrayOfStrings"`
FieldAnonymousStruct StructOptionalFieldsSomeStructFieldAnonymousStruct `tfsdk:"FieldAnonymousStruct"`
 }

type SomeOtherStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type StructOptionalFieldsSomeStructFieldAnonymousStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }



var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_ref": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},

"field_string": schema.StringAttribute{
 Optional: true,
},

"operator": schema.StringAttribute{
 Required: true,
},

"field_array_of_strings": schema.ListAttribute{
 ElementType: types.StringType,
},

"field_anonymous_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},

},
},
"some_other_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
"struct_optional_fields_some_struct_field_anonymous_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
}
package struct_complex_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// This struct does things.
type SomeStruct struct {
 FieldRef SomeOtherStruct `tfsdk:"FieldRef"`
FieldDisjunctionOfScalars StringOrBool `tfsdk:"FieldDisjunctionOfScalars"`
FieldMixedDisjunction StringOrSomeOtherStruct `tfsdk:"FieldMixedDisjunction"`
FieldDisjunctionWithNull types.String `tfsdk:"FieldDisjunctionWithNull"`
Operator types.String `tfsdk:"Operator"`
FieldArrayOfStrings types.List `tfsdk:"FieldArrayOfStrings"`
FieldMapOfStringToString types.Map `tfsdk:"FieldMapOfStringToString"`
FieldAnonymousStruct StructComplexFieldsSomeStructFieldAnonymousStruct `tfsdk:"FieldAnonymousStruct"`
FieldRefToConstant ConnectionPath `tfsdk:"fieldRefToConstant"`
 }

const ConnectionPath = "straight"

type SomeOtherStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type StructComplexFieldsSomeStructFieldAnonymousStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }



type StringOrBool struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
 }

type StringOrSomeOtherStruct struct {
 String types.String `tfsdk:"String"`
SomeOtherStruct SomeOtherStruct `tfsdk:"SomeOtherStruct"`
 }

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Description: `
This struct does things.
`,
Attributes: map[string]schema.Attribute{
"field_ref": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},

"field_disjunction_of_scalars": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"string": schema.StringAttribute{
 Optional: true,
},

"bool": schema.BoolAttribute{
 Optional: true,
},

},
},

"field_mixed_disjunction": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"string": schema.StringAttribute{
 Optional: true,
},

"some_other_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},

},
},

"field_disjunction_with_null": schema.StringAttribute{
 Optional: true,
},

"operator": schema.StringAttribute{
 Required: true,
},

"field_array_of_strings": schema.ListAttribute{
 ElementType: types.StringType,
},

"field_map_of_string_to_string": schema.MapAttribute{
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

"field_ref_to_constant": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("straight"),
},

},
},
"connection_path": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("straight"),
},
"some_other_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
"struct_complex_fields_some_struct_field_anonymous_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
}
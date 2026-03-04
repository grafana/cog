package disjunctions

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// Refresh rate or disabled.
type RefreshRate = StringOrBool

type StringOrNull types.String

type SomeStruct struct {
 Type types.String `tfsdk:"Type"`
FieldAny types.Object `tfsdk:"FieldAny"`
 }

type BoolOrRef = BoolOrSomeStruct

type SomeOtherStruct struct {
 Type types.String `tfsdk:"Type"`
Foo types.String `tfsdk:"Foo"`
 }

type YetAnotherStruct struct {
 Type types.String `tfsdk:"Type"`
Bar types.Number `tfsdk:"Bar"`
 }

type SeveralRefs = SomeStructOrSomeOtherStructOrYetAnotherStruct

type StringOrBool struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
 }

type BoolOrSomeStruct struct {
 Bool types.Bool `tfsdk:"Bool"`
SomeStruct SomeStruct `tfsdk:"SomeStruct"`
 }

type SomeStructOrSomeOtherStructOrYetAnotherStruct struct {
 SomeStruct SomeStruct `tfsdk:"SomeStruct"`
SomeOtherStruct SomeOtherStruct `tfsdk:"SomeOtherStruct"`
YetAnotherStruct YetAnotherStruct `tfsdk:"YetAnotherStruct"`
 }

var SpecAttributes = map[string]schema.Attribute{
"string_or_null": schema.StringAttribute{
 Optional: true,
},
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("some-struct"),
},

"field_any": schema.ObjectAttribute{
 Required: true,
},

},
},
"some_other_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("some-other-struct"),
},

"foo": schema.StringAttribute{
 Required: true,
},

},
},
"yet_another_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("yet-another-struct"),
},

"bar": schema.NumberAttribute{
 Required: true,
},

},
},
}
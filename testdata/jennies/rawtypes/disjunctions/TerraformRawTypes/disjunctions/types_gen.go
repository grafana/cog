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

var SomeStructAttributes = map[string]schema.Attribute{
"type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("some-struct"),
},

"field_any": schema.ObjectAttribute{
 Required: true,
},

}

var SomeOtherStructAttributes = map[string]schema.Attribute{
"type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("some-other-struct"),
},

"foo": schema.StringAttribute{
 Required: true,
},

}

var YetAnotherStructAttributes = map[string]schema.Attribute{
"type": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("yet-another-struct"),
},

"bar": schema.NumberAttribute{
 Required: true,
},

}

var StringOrBoolAttributes = map[string]schema.Attribute{
"string": schema.StringAttribute{
 Optional: true,
},

"bool": schema.BoolAttribute{
 Optional: true,
},

}

var BoolOrSomeStructAttributes = map[string]schema.Attribute{
"bool": schema.BoolAttribute{
 Optional: true,
},

"some_struct": schema.SingleNestedAttribute{
Optional: true,
Attributes: SomeStructAttributes,
},

}

var SomeStructOrSomeOtherStructOrYetAnotherStructAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Optional: true,
Attributes: SomeStructAttributes,
},

"some_other_struct": schema.SingleNestedAttribute{
Optional: true,
Attributes: SomeOtherStructAttributes,
},

"yet_another_struct": schema.SingleNestedAttribute{
Optional: true,
Attributes: YetAnotherStructAttributes,
},

}

var SpecAttributes = map[string]schema.Attribute{
"string_or_null": schema.StringAttribute{
 Optional: true,
},
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: SomeStructAttributes,
},
"some_other_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: SomeOtherStructAttributes,
},
"yet_another_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: YetAnotherStructAttributes,
},
}
package disjunctions_of_scalars_and_refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

type DisjunctionOfScalarsAndRefs = StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB

type MyRefA struct {
 Foo types.String `tfsdk:"foo"`
 }

type MyRefB struct {
 Bar types.Int64 `tfsdk:"bar"`
 }

type StringOrBoolOrArrayOfStringOrMyRefAOrMyRefB struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
ArrayOfString types.List `tfsdk:"ArrayOfString"`
MyRefA MyRefA `tfsdk:"MyRefA"`
MyRefB MyRefB `tfsdk:"MyRefB"`
 }

var MyRefAAttributes = map[string]schema.Attribute{
"foo": schema.StringAttribute{
 Required: true,
},

}

var MyRefBAttributes = map[string]schema.Attribute{
"bar": schema.Int64Attribute{
 Required: true,
},

}

var StringOrBoolOrArrayOfStringOrMyRefAOrMyRefBAttributes = map[string]schema.Attribute{
"string": schema.StringAttribute{
 Optional: true,
Default: stringdefault.StaticString("a"),
},

"bool": schema.BoolAttribute{
 Optional: true,
},

"array_of_string": schema.ListAttribute{
 ElementType: types.StringType,
},

"my_ref_a": schema.SingleNestedAttribute{
Optional: true,
Attributes: MyRefAAttributes,
},

"my_ref_b": schema.SingleNestedAttribute{
Optional: true,
Attributes: MyRefBAttributes,
},

}

var SpecAttributes = map[string]schema.Attribute{
"my_ref_a": schema.SingleNestedAttribute{
Required: true,
Attributes: MyRefAAttributes,
},
"my_ref_b": schema.SingleNestedAttribute{
Required: true,
Attributes: MyRefBAttributes,
},
}
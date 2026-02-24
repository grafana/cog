package disjunctions_of_scalars_and_refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


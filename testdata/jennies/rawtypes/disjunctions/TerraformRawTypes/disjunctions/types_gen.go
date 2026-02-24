package disjunctions

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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


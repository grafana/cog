package withdashes

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

// Refresh rate or disabled.
type RefreshRate = StringOrBool

type StringOrBool struct {
 String types.String `tfsdk:"String"`
Bool types.Bool `tfsdk:"Bool"`
 }


package constant_reference_as_default

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

const ConstantRefString = "AString"

type MyStruct struct {
 AString types.String `tfsdk:"aString"`
OptString types.String `tfsdk:"optString"`
 }


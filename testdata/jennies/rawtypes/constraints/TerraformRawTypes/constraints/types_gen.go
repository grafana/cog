package constraints

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type SomeStruct struct {
 Id types.Int64 `tfsdk:"id"`
MaybeId types.Int64 `tfsdk:"maybeId"`
Title types.String `tfsdk:"title"`
RefStruct RefStruct `tfsdk:"refStruct"`
 }

type RefStruct struct {
 Labels types.Map `tfsdk:"labels"`
Tags types.List `tfsdk:"tags"`
 }


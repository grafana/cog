package variant_dataquery

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type Query struct {
 Expr types.String `tfsdk:"expr"`
Instant types.Bool `tfsdk:"instant"`
 }


package time_hint

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type ObjTime types.String

type ObjWithTimeField struct {
 RegisteredAt types.String `tfsdk:"registeredAt"`
 }


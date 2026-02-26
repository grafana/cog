package time_hint

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	timetypes "/github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

type ObjTime types.String

type ObjWithTimeField struct {
 RegisteredAt timetypes.RFC3339 `tfsdk:"registeredAt"`
 }


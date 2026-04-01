package time_hint

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	timetypes "github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type ObjTime types.String

type ObjWithTimeField struct {
 RegisteredAt timetypes.RFC3339 `tfsdk:"registeredAt"`
Duration timetypes.GoDurationType `tfsdk:"duration"`
 }

var SpecAttributes = map[string]schema.Attribute{
"obj_time": schema.StringAttribute{
 Required: true,
CustomType: timetypes.RFC3339Type{},
},
"obj_with_time_field": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"registered_at": schema.StringAttribute{
 Required: true,
CustomType: timetypes.RFC3339Type{},
},

"duration": schema.StringAttribute{
 Required: true,
CustomType: timetypes.GoDurationType{},
},

},
},
}
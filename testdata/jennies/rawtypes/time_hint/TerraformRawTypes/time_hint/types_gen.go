package time_hint

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	timetypes "/github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
)

type ObjTime types.String

type ObjWithTimeField struct {
 RegisteredAt timetypes.RFC3339 `tfsdk:"registeredAt"`
 }

var SpecAttributes = map[string]schema.Attribute{
"objtime": schema.StringAttribute{
 Required: true
 
}"objwithtimefield": types.ObjectAttributes{
Required: true,
AttributeTypes: map[string]attr.Type{
"registeredAt": types.StringType,
},
}
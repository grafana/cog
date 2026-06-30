package time_hint

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	timetypes "github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type ObjTimeModel = types.String
var ObjTimeType = timetypes.RFC3339


type ObjWithTimeFieldModel struct {
	RegisteredAt timetypes.RFC3339 `tfsdk:"registered_at"`
	Duration timetypes.GoDurationType `tfsdk:"duration"`
}
var ObjWithTimeFieldType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"registered_at": timetypes.RFC3339,
		"duration": timetypes.GoDurationType,
	},
}
var ObjWithTimeFieldSchema = schema.SingleNestedBlock{
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
	Blocks: map[string]schema.Block{
},
}


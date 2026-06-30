package intersections

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	booldefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

type IntersectionsModel unknown
var IntersectionsType = unknown


type SomeStructModel struct {
	FieldBool types.Bool `tfsdk:"field_bool"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field_bool": types.BoolType,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field_bool": schema.BoolAttribute{
	Optional: true,
	Default: booldefault.StaticBool(true),
	Computed: true,
},
},
	Blocks: map[string]schema.Block{
},
}

// Base properties for all metrics
type CommonModel struct {
	// The metric name
	Name types.String `tfsdk:"name"`
	// The metric type
	Type types.String `tfsdk:"type"`
	// The type of data the metric contains
	Contains types.String `tfsdk:"contains"`
}
var CommonType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name": types.StringType,
		"type": types.StringType,
		"contains": types.StringType,
	},
}
var CommonSchema = schema.SingleNestedBlock{
	Description: "Base properties for all metrics",
	MarkdownDescription: "Base properties for all metrics",
	Attributes: map[string]schema.Attribute{
	"name": schema.StringAttribute{
	Required: true,
	Description: "The metric name",
	MarkdownDescription: "The metric name",
},
	"type": schema.StringAttribute{
	Required: true,
	Validators: []validator.String{
stringvalidator.OneOf("counter", "gauge"),
},
	Description: "The metric type",
	MarkdownDescription: "The metric type",
},
	"contains": schema.StringAttribute{
	Required: true,
	Validators: []validator.String{
stringvalidator.OneOf("default", "time"),
},
	Description: "The type of data the metric contains",
	MarkdownDescription: "The type of data the metric contains",
},
},
	Blocks: map[string]schema.Block{
},
}

// Counter metric combining common properties with specific values
type CounterModel unknown
var CounterType = unknown


type CommonTypeModel types.String
var CommonTypeType = types.StringType


type CommonContainsModel types.String
var CommonContainsType = types.StringType



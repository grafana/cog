package intersections

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	booldefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)



type SomeStruct struct {
 FieldBool types.Bool `tfsdk:"fieldBool"`
 }

// Base properties for all metrics
type Common struct {
 // The metric name
Name types.String `tfsdk:"name"`
// The metric type
Type types.String `tfsdk:"type"`
// The type of data the metric contains
Contains types.String `tfsdk:"contains"`
 }

// Counter metric combining common properties with specific values






var SomeStructAttributes = map[string]schema.Attribute{
"field_bool": schema.BoolAttribute{
 Required: true,
Default: booldefault.StaticBool(true),
},

}

var CommonAttributes = map[string]schema.Attribute{
"name": schema.StringAttribute{
 Required: true,
},

"type": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("counter", "gauge"),
},

},

"contains": schema.StringAttribute{
 Required: true,
Validators: []validator.String{
stringvalidator.OneOf("default", "time"),
},

},

}

var SpecAttributes = map[string]schema.Attribute{
"some_struct": schema.SingleNestedAttribute{
Required: true,
Attributes: SomeStructAttributes,
},
"common": schema.SingleNestedAttribute{
Required: true,
Description: `
Base properties for all metrics
`,
Attributes: CommonAttributes,
},
}
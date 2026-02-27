package intersections

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "/github.com/hashicorp/terraform-plugin-framework/resource/schema"
	attr "/github.com/hashicorp/terraform-plugin-framework/attr"
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






var SpecAttributes = map[string]schema.Attribute{
"someStruct": schema.ObjectAttribute{
Required: true,
AttributeTypes: map[string]attr.Type{
"fieldBool": types.BoolType,
},
},
"common": schema.ObjectAttribute{
Required: true,
Description: `
Base properties for all metrics
`,
AttributeTypes: map[string]attr.Type{
"name": types.StringType,
"type": types.StringType,
"contains": types.StringType,
},
},
}
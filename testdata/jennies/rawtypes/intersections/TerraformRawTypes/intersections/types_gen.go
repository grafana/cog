package intersections

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
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







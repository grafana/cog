package self_referential_struct

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Node represents a node in a singly-linked list.
// The next field points to the following node, or is absent if this is the last node.
type Node struct {
 Value types.String `tfsdk:"value"`
Next Node `tfsdk:"next"`
 }

var NodeAttributes = map[string]schema.Attribute{
"value": schema.StringAttribute{
 Required: true,
},

"next": schema.SingleNestedAttribute{
Optional: true,
Attributes: map[string]schema.Attribute{
"value": schema.StringAttribute{
 Required: true,
},

},
},

}

var SpecAttributes = map[string]schema.Attribute{
"node": schema.SingleNestedAttribute{
Required: true,
Description: `
Node represents a node in a singly-linked list.
The next field points to the following node, or is absent if this is the last node.
`,
Attributes: NodeAttributes,
},
}

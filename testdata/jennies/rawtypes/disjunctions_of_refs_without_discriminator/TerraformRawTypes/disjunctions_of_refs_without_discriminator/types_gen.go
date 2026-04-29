package disjunctions_of_refs_without_discriminator

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type DisjunctionWithoutDiscriminator types.Object

type TypeA struct {
 FieldA types.String `tfsdk:"fieldA"`
 }

type TypeB struct {
 FieldB types.Int64 `tfsdk:"fieldB"`
 }

var SpecAttributes = map[string]schema.Attribute{
"disjunction_without_discriminator": schema.ObjectAttribute{
 Required: true,
},
"type_a": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_a": schema.StringAttribute{
 Required: true,
},

},
},
"type_b": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field_b": schema.Int64Attribute{
 Required: true,
},

},
},
}
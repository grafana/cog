package nullable_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

type Struct struct {
 A MyObject `tfsdk:"a"`
B MyObject `tfsdk:"b"`
C types.String `tfsdk:"c"`
D types.List `tfsdk:"d"`
E types.Map `tfsdk:"e"`
F NullableFieldsStructF `tfsdk:"f"`
G types.String `tfsdk:"g"`
 }

type MyObject struct {
 Field types.String `tfsdk:"field"`
 }

const ConstantRef = "hey"

type NullableFieldsStructF struct {
 A types.String `tfsdk:"a"`
 }

var SpecAttributes = map[string]schema.Attribute{
"struct": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"a": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field": schema.StringAttribute{
 Required: true,
},

},
},

"b": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field": schema.StringAttribute{
 Required: true,
},

},
},

"c": schema.StringAttribute{
 Optional: true,
},

"d": schema.ListAttribute{
 ElementType: types.StringType,
},

"e": schema.MapAttribute{
 ElementType: types.StringType,
},

"f": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"a": schema.StringAttribute{
 Required: true,
},

},
},

"g": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("hey"),
},

},
},
"my_object": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"field": schema.StringAttribute{
 Required: true,
},

},
},
"constant_ref": schema.StringAttribute{
 Required: true,
Default: stringdefault.StaticString("hey"),
},
"nullable_fields_struct_f": schema.SingleNestedAttribute{
Required: true,
Attributes: map[string]schema.Attribute{
"a": schema.StringAttribute{
 Required: true,
},

},
},
}
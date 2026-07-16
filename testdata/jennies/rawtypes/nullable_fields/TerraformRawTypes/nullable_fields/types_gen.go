package nullable_fields

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type StructModel struct {
	A types.Object `tfsdk:"a"`
	B types.Object `tfsdk:"b"`
	C types.String `tfsdk:"c"`
	D types.List `tfsdk:"d"`
	E types.Map `tfsdk:"e"`
	F types.Object `tfsdk:"f"`
}
var StructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"a": MyObjectType,
		"b": MyObjectType,
		"c": types.StringType,
		"d": types.ListType{
	ElemType: types.StringType,
},
		"e": types.MapType{
	ElemType: types.StringType,
},
		"f": NullableFieldsStructFType,
	},
}
var StructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"c": schema.StringAttribute{
	Optional: true,
},
	"d": schema.ListAttribute{
	ElementType: types.StringType,
	Optional: true,
},
	"e": schema.MapAttribute{
	ElementType: types.StringType,
	Required: true,
},
},
	Blocks: map[string]schema.Block{
	"a": schema.SingleNestedBlock{
		Attributes: MyObjectSchema.Attributes,
		Blocks: MyObjectSchema.Blocks,
		Validators: MyObjectSchema.Validators,
	},
	"b": schema.SingleNestedBlock{
		Attributes: MyObjectSchema.Attributes,
		Blocks: MyObjectSchema.Blocks,
		Validators: MyObjectSchema.Validators,
	},
	"f": schema.SingleNestedBlock{
		Attributes: NullableFieldsStructFSchema.Attributes,
		Blocks: NullableFieldsStructFSchema.Blocks,
		Validators: NullableFieldsStructFSchema.Validators,
	},
},
}

type MyObjectModel struct {
	Field types.String `tfsdk:"field"`
}
var MyObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"field": types.StringType,
	},
}
var MyObjectSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"field": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}

const ConstantRefModel = "hey"
var ConstantRefType = types.StringType


type NullableFieldsStructFModel struct {
	A types.String `tfsdk:"a"`
}
var NullableFieldsStructFType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"a": types.StringType,
	},
}
var NullableFieldsStructFSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"a": schema.StringAttribute{
	Required: true,
},
},
	Blocks: map[string]schema.Block{
},
}


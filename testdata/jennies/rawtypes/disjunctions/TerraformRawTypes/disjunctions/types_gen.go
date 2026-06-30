package disjunctions

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
	attr "github.com/hashicorp/terraform-plugin-framework/attr"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	stringdefault "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	validator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Refresh rate or disabled.
type RefreshRateModel = types.Object
var RefreshRateType = StringOrBoolType


type StringOrNullModel = types.String
var StringOrNullType = types.StringType


type SomeStructModel struct {
	Type types.String `tfsdk:"type"`
	FieldAny types.String `tfsdk:"field_any"`
}
var SomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"field_any": types.StringType,
	},
}
var SomeStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"type": schema.StringAttribute{
	Optional: true,
	Default: stringdefault.StaticString("some-struct"),
	Computed: true,
},
	"field_any": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type BoolOrRefModel = types.Object
var BoolOrRefType = BoolOrSomeStructType


type SomeOtherStructModel struct {
	Type types.String `tfsdk:"type"`
	Foo types.String `tfsdk:"foo"`
}
var SomeOtherStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"foo": types.StringType,
	},
}
var SomeOtherStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"type": schema.StringAttribute{
	Optional: true,
	Default: stringdefault.StaticString("some-other-struct"),
	Computed: true,
},
	"foo": schema.StringAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type YetAnotherStructModel struct {
	Type types.String `tfsdk:"type"`
	Bar types.Number `tfsdk:"bar"`
}
var YetAnotherStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		"bar": types.NumberType,
	},
}
var YetAnotherStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"type": schema.StringAttribute{
	Optional: true,
	Default: stringdefault.StaticString("yet-another-struct"),
	Computed: true,
},
	"bar": schema.NumberAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
}

type SeveralRefsModel = types.Object
var SeveralRefsType = SomeStructOrSomeOtherStructOrYetAnotherStructType


type StringOrBoolModel struct {
	String types.String `tfsdk:"string"`
	Bool types.Bool `tfsdk:"bool"`
}
var StringOrBoolType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"string": types.StringType,
		"bool": types.BoolType,
	},
}
var StringOrBoolSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"string": schema.StringAttribute{
	Optional: true,
},
	"bool": schema.BoolAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
},
	Validators: []validator.Object{
		(1, "string", "bool"),
	},
}

type BoolOrSomeStructModel struct {
	Bool types.Bool `tfsdk:"bool"`
	SomeStruct types.Object `tfsdk:"some_struct"`
}
var BoolOrSomeStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"bool": types.BoolType,
		"some_struct": SomeStructType,
	},
}
var BoolOrSomeStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
	"bool": schema.BoolAttribute{
	Optional: true,
},
},
	Blocks: map[string]schema.Block{
	"some_struct": schema.SingleNestedBlock{
		Attributes: SomeStructSchema.Attributes,
		Blocks: SomeStructSchema.Blocks,
		Validators: SomeStructSchema.Validators,
	},
},
}

type SomeStructOrSomeOtherStructOrYetAnotherStructModel struct {
	SomeStruct types.Object `tfsdk:"some_struct"`
	SomeOtherStruct types.Object `tfsdk:"some_other_struct"`
	YetAnotherStruct types.Object `tfsdk:"yet_another_struct"`
}
var SomeStructOrSomeOtherStructOrYetAnotherStructType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"some_struct": SomeStructType,
		"some_other_struct": SomeOtherStructType,
		"yet_another_struct": YetAnotherStructType,
	},
}
var SomeStructOrSomeOtherStructOrYetAnotherStructSchema = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
},
	Blocks: map[string]schema.Block{
	"some_struct": schema.SingleNestedBlock{
		Attributes: SomeStructSchema.Attributes,
		Blocks: SomeStructSchema.Blocks,
		Validators: []validator.Object{
			("type", "field_any"),
		},
	},
	"some_other_struct": schema.SingleNestedBlock{
		Attributes: SomeOtherStructSchema.Attributes,
		Blocks: SomeOtherStructSchema.Blocks,
		Validators: []validator.Object{
			("type", "foo"),
		},
	},
	"yet_another_struct": schema.SingleNestedBlock{
		Attributes: YetAnotherStructSchema.Attributes,
		Blocks: YetAnotherStructSchema.Blocks,
		Validators: []validator.Object{
			("type", "bar"),
		},
	},
},
	Validators: []validator.Object{
		(1, "some_struct", "some_other_struct", "yet_another_struct"),
	},
}


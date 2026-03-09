package constraints

import (
	int64validator "github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	stringvalidator "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SomeStruct struct {
	Id        types.Int64  `tfsdk:"id"`
	MaybeId   types.Int64  `tfsdk:"maybeId"`
	Title     types.String `tfsdk:"title"`
	RefStruct RefStruct    `tfsdk:"refStruct"`
}

type RefStruct struct {
	Labels types.Map  `tfsdk:"labels"`
	Tags   types.List `tfsdk:"tags"`
}

var SpecAttributes = map[string]schema.Attribute{
	"some_struct": schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
				[]validator.Int64{
					int64validator.AtLeast(5),
					int64validator.AtMost(9),
				},
			},

			"maybe_id": schema.Int64Attribute{
				Optional: true,
				[]validator.Int64{
					int64validator.AtLeast(5),
					int64validator.AtMost(9),
				},
			},

			"title": schema.StringAttribute{
				Required: true,
				[]validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"ref_struct": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"labels": schema.MapAttribute{
						ElementType: types.StringType,
					},

					"tags": schema.ListAttribute{
						ElementType: types.StringType,
					},
				},
			},
		},
	},
	"ref_struct": schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"labels": schema.MapAttribute{
				ElementType: types.StringType,
			},

			"tags": schema.ListAttribute{
				ElementType: types.StringType,
			},
		},
	},
}

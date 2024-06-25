package dashboard

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type DashboardDataSourceModel struct {
     Title types.String `tfsdk:"title"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type DataSourceRefDataSourceModel struct {
     Type types.String `tfsdk:"type"`
     Uid types.String `tfsdk:"uid"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type FieldConfigSourceDataSourceModel struct {
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type FieldConfigDataSourceModel struct {
     Unit types.String `tfsdk:"unit"`
     Custom types.Object `tfsdk:"custom"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type PanelDataSourceModel struct {
     Title types.String `tfsdk:"title"`
     Type types.String `tfsdk:"type"`
     Options types.Object `tfsdk:"options"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


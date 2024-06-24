package variant_dataquery

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type QueryDataSourceModel struct {
     Expr types.String `tfsdk:"expr"`
     Instant types.Bool `tfsdk:"instant"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


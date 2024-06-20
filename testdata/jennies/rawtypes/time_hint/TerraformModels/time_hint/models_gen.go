package time_hint

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type ObjWithTimeFieldDataSourceModel struct {
     RegisteredAt types.String `tfsdk:"registered_at"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


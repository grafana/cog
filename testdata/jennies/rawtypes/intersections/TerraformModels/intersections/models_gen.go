package intersections

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldBool types.Bool `tfsdk:"field_bool"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


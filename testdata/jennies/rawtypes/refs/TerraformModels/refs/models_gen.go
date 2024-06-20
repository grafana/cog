package refs

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldAny types.Object `tfsdk:"field_any"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


package withdashes

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldAny types.Object `tfsdk:"field_any"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type StringOrBoolDataSourceModel struct {
     String types.String `tfsdk:"string"`
     Bool types.Bool `tfsdk:"bool"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


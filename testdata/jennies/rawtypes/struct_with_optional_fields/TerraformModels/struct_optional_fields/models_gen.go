package struct_optional_fields

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldString types.String `tfsdk:"field_string"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type SomeOtherStructDataSourceModel struct {
     FieldAny types.Object `tfsdk:"field_any"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


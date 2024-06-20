package defaults

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type NestedStructDataSourceModel struct {
     StringVal types.String `tfsdk:"string_val"`
     IntVal types.Int64 `tfsdk:"int_val"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type StructDataSourceModel struct {
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


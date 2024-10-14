package constraints

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     Id types.Int64 `tfsdk:"id"`
     MaybeId types.Int64 `tfsdk:"maybe_id"`
     Title types.String `tfsdk:"title"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type RefStructDataSourceModel struct {
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


package struct_complex_fields

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     FieldDisjunctionWithNull types.String `tfsdk:"field_disjunction_with_null"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type SomeOtherStructDataSourceModel struct {
     FieldAny types.Object `tfsdk:"field_any"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type StringOrBoolDataSourceModel struct {
     String types.String `tfsdk:"string"`
     Bool types.Bool `tfsdk:"bool"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type StringOrSomeOtherStructDataSourceModel struct {
     String types.String `tfsdk:"string"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


package disjunctions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type SomeStructDataSourceModel struct {
     Type types.String `tfsdk:"type"`
     FieldAny types.Object `tfsdk:"field_any"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type SomeOtherStructDataSourceModel struct {
     Type types.String `tfsdk:"type"`
     Foo types.String `tfsdk:"foo"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type YetAnotherStructDataSourceModel struct {
     Type types.String `tfsdk:"type"`
     Bar types.Int64 `tfsdk:"bar"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type StringOrBoolDataSourceModel struct {
     String types.String `tfsdk:"string"`
     Bool types.Bool `tfsdk:"bool"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type BoolOrSomeStructDataSourceModel struct {
     Bool types.Bool `tfsdk:"bool"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}

type SomeStructOrSomeOtherStructOrYetAnotherStructDataSourceModel struct {
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


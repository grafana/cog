package variant_panelcfg_only_options

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)


type OptionsDataSourceModel struct {
     Content types.String `tfsdk:"content"`
  TemporaryScalarPlaceholder types.Bool // @TODO Remove this once non-scalars are implemented
}


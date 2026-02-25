package refs

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type RefToSomeStruct = SomeStruct

type RefToSomeStructFromOtherPackage = unknown


package arrays

import (
	 "github.com/hashicorp/terraform-plugin-framework/types"
)

// List of tags, maybe?
type ArrayOfStrings types.List

type SomeStruct struct {
 FieldAny types.Object `tfsdk:"FieldAny"`
 }

type ArrayOfRefs []SomeStruct

type ArrayOfArrayOfNumbers types.List


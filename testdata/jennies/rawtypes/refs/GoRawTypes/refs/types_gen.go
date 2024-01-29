package refs

import (
	otherpkg "github.com/grafana/cog/generated/otherpkg"
)

type SomeStruct struct {
	FieldAny any `json:"FieldAny"`
}

type RefToSomeStruct SomeStruct

type RefToSomeStructFromOtherPackage otherpkg.SomeDistantStruct


package constraints

import "strings"

SomeStruct: {
	id: int64 & >= 5 & <10
	title: strings.MinRunes(1) & {
		string
	}
}

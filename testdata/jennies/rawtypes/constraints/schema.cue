package constraints

import "strings"

SomeStruct: {
	id: int64 & >= 5 & <10
	maybeId?: int64 & >= 5 & <10
	greaterThanZero: int64 & >= 0 & <3
	negative: int64 & >=-19 & <10
	title: strings.MinRunes(1) & string
	labels: [string]: (string & strings.MinRunes(1))
	tags: [...(string & strings.MinRunes(1))]
}

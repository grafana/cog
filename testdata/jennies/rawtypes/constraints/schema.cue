package constraints

import (
	"strings"
	"list"
)

SomeStruct: {
	id: int64 & >= 5 & <10
	maybeId?: int64 & >= 5 & <10
	greaterThanZero: int64 & >= 0 & <3
	negative: int64 & >=-19 & <10
	title: strings.MinRunes(1) & string
	labels: [string]: (string & strings.MinRunes(1))
	tags: [...(string & strings.MinRunes(1))]
	regex: string & =~ "^[a-zA-Z0-9_-]+$"
	negativeRegex: string & !~ "^[a-zA-Z0-9_-]+$"
	minMaxList: list.MinItems(1) & list.MaxItems(64) & [...string]
	uniqueList: list.UniqueItems() & [...string]
	fullConstraintList: list.MinItems(2) & list.MaxItems(10) & list.UniqueItems() & [...int64]
}

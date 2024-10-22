package constraints

import "strings"

SomeStruct: {
	id: int64 & >= 5 & <10
	maybeId?: int64 & >= 5 & <10
	title: strings.MinRunes(1) & {
		string
	}
	refStruct?: refStruct
}


refStruct: {
	labels: [string]: (string & strings.MinRunes(1))
	tags: [...(string & strings.MinRunes(1))]
}

package validation

import "strings"

Dashboard: {
	uid?: string & strings.MinRunes(1)
	id?: int & >0
	title: string & strings.MinRunes(1)
	tags: [...(string & strings.MinRunes(1))]
	labels: [string]: (string & strings.MinRunes(1))
	panels: [...#Panel]
}

#Panel: {
	title: string & strings.MinRunes(1)
}

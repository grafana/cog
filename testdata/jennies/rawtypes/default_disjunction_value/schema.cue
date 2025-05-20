package default_disjunction_value

DisjunctionClasses: *ValueA | ValueB | ValueC
DisjunctionConstants: "abc" | *1 | true

ValueA: {
	type: "A"
	anArray: [...string]
	otherRef: ValueB
}

ValueB: {
	type: "B"
	aMap: [string]: int64
	def: *1 | "a" | bool
}

ValueC: {
	type: "C"
	other: float32
}

package disjunctions_of_refs

DisjunctionOfRefs: MyRefA | MyRefB

MyRefA: {
	type: "A"
	foo: string
}

MyRefB: {
	type: "B"
	bar: int64
}

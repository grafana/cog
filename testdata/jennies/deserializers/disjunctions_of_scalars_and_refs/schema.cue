package disjunctions_of_scalars_and_refs

DisjunctionOfScalarsAndRefs: "a" | bool | [...string] | MyRefA | MyRefB | {}

MyRefA: {
	foo: string
}

MyRefB: {
	bar: int64
}

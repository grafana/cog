package disjunction_anonymous

MyStruct: {
	scalars: string | bool | number | int
	sameKind: "a" | "b" | "c"
	refs: StructA | StructB
	mixed: StructA | string | int
}

StructA: {
	field: string
}

StructB: {
	type: int
}

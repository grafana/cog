package nullable_fields

Struct: {
	a: MyObject | null
	b?: MyObject | null
	c: string | null
	d: [...string] | null
	e: [string]: string | null
	f: {
		a: string
	} | null
	g: ConstantRef | null
}

ConstantRef: "hey"

MyObject: {
	field: string
}

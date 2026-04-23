package sandbox

StringEnum: "a" | "b" | "c"
StringEnumWithDefault: *"a" | "b" | "c"

SomeStruct: {
	data: [StringEnum]: string
}

SomeStructWithDefaultEnum: {
	data: [StringEnumWithDefault]: string
}

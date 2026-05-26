package enums_as_map_index

#StringEnum: "a" | "b" | "c"
#StringEnumWithDefault: *"a" | "b" | "c"

#SomeStruct: {
	data: [#StringEnum]: string
}

#SomeStructWithDefaultEnum: {
	data: [#StringEnumWithDefault]: string
}

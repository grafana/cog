package constant_references

Enum: "ValueA" | "ValueB" | "ValueC"

ParentStruct: {
	myEnum: Enum
}

Struct: {
	ParentStruct
	myValue: string
}

StructA: {
	ParentStruct
	myEnum: Enum & "ValueA"
}

StructB: {
	Struct
	myEnum: Enum & "ValueB"
}

package reference_of_reference

MyStruct: {
	field?: OtherStruct
}

OtherStruct: AnotherStruct

AnotherStruct: {
	a: string
}

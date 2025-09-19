package constant_reference_as_default

ConstantRefString: "AString"

MyStruct: {
	aString: ConstantRefString
	optString?: ConstantRefString
}

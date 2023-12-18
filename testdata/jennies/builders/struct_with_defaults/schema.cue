package struct_with_defaults

#NestedStruct: {
	stringVal: string
	intVal: int64
}

#Struct: {
	allFields: #NestedStruct | *{
  	stringVal: "hello"
  	intVal: 3
  }
  partialFields: #NestedStruct | *{
  	intVal: 4
  }
  emptyFields: #NestedStruct | *{}

  complexField: {
		uid: string
		nested: {
			nestedVal: string
		}
		array: [...string]
	} | *{ uid: "myUID", nested: { nestedVal: "nested"}, array: ["hello"]}

	partialComplexField: {
		uid: string
		aVal: int64
	} | *{ uid: "myUID" }
}

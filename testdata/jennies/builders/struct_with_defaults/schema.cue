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
}

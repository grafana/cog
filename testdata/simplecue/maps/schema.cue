intStringMap: {
    [string]: int
}
stringStringMap: {
    [string]: string
}
#foo: {
    bar: string
		stringToAny: {
				[string]: _
    }
}
stringRefMap: [string]: #foo
stringToMapOfMap: [string]: {[string]: bool}

incompleteObjectIsNotAMap: {
	foo: string
	...
}

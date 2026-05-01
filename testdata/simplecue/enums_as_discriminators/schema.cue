package enums_as_discriminators

#Enum1: "A" | "B"
#Enum2: *"C" | "D"

Struct: {
	aOrB: MyStructA | MyStructB
	cOrD: MyStructC | MyStructD
	disjAorB: DisjAB
	disjCorD: DisjCD
}

MyStructA: {
	type: #Enum1 & "A"
	field: string
}

MyStructB: {
	type: #Enum1 & "B"
	something: int
}

MyStructC: {
	type: #Enum2 & "C"
	field: string
}

MyStructD: {
	type: #Enum2 & "D"
	something: int
}

DisjAB: MyStructA | MyStructB
DisjCD: MyStructC | MyStructD

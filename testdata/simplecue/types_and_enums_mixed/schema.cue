stringEnum: "a" | "b" | "c"
expressionTypes: expressionA | expressionB | expressionC @cog(kind="type")

container: {
	value: expressionTypes
	disj: expressionA | expressionC
}

expressionA: {
	enumA: stringEnum & "a"
}

expressionB: {
	enumB: stringEnum & "b"
}

expressionC: {
	enumC: stringEnum & "c"
}



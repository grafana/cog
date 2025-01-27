stringEnum: "a" | "b" | "c"

container: {
	stringA: stringEnum
}

withoutOverride: {
	container
	stringA: stringEnum & "a"
}



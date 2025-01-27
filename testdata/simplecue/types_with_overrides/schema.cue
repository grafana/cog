stringType: "a" | "b" | "c"  @cog(kind="type")

container: {
	stringA: stringType
}

withoutOverride: {
	container
	stringA: stringType & "a"
}



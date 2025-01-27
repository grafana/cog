stringEnum: "a" | "b" | "c"
intEnum: 1 | 2 | 3 @cog(kind="enum", memberNames="X|Y|Z")

container: {
	enumA: stringEnum & "a"
	enumB: stringEnum & "b"
	enumC: stringEnum & "c"
	enum1: intEnum & 1
	enum2: intEnum & 2
	enum3: intEnum & 3
}



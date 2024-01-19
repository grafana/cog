#IntEnum: 0 | 1 | *2 @cog(kind="enum",memberNames="firstValue|secondValue|thirdValue")

container: {
    StringEnum: "foo" | "bar" | "baz"
    TheIntEnum: #IntEnum
}

anything: { ... }
foo: {
    arrayOfAnything: [...]
}

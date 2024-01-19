StringEnum: "a" | "b" | "c"
IntEnum: 1 | 2 | 3 @cog(kind="enum", memberNames="one|two|three")

container: {
    from: string | *"now-6h"
    editable: bool | *true
    IntEnum: 0 | 1 | *2 @cog(kind="enum",memberNames="firstValue|secondValue|thirdValue")
    Number: int64 | *5
    repeatDirection: *"h" | "v"
    tags: [...string] | *["default", "tags"]
    stringEnum: StringEnum & (*"c" | _)
    intEnum: IntEnum & (*2 | _)
}

#refStruct: {
    things: string
}

struct: {
    things: string
}

container: {
    ints: [...int]
    strings: [...string]
    refs: [...#refStruct]
    structs: [...struct]
    arrayOfAnything: [...]
}

package enums_as_map_index

#StringEnum: "a" | "b" | "c"

#SomeStruct: {
    field: string
}

EnumIndexedMap: [#StringEnum]: string
EnumIndexedMapOfStructs: [#StringEnum]: #SomeStruct

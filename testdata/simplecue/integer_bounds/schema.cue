#PositiveInt32: int32 & >0

container: {
    // Bounds keep the declared width, and are all kept as constraints.
    intWithMaxInt32: int & >0 & <=2147483647
    int64WithMaxInt32: int64 & >0 & <=2147483647
    intWithSmallRange: int & >0 & <100
    intWithInclusiveRange: int & >=0 & <=2147483647
    intWithDefault: int & >0 & <=2147483647 | *30
    int64WithDefault: int64 & >=1 & <=2147483647 | *30

    // Declared fixed-width types keep their width, without restating its limits.
    int32Positive: int32 & >0
    int32Range: int32 & >0 & <100
    uint32WithDefault: uint32 & >0 | *9
    int32Ref: #PositiveInt32
    int32RefWithBound: #PositiveInt32 & <=1000
    int8OrInt16: int8 | int16
    uintWithMax: uint & <100
}

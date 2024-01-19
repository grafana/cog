container: {
    constantInt: 42
    constantFloat: 42.24
    rowHeight: float & >=0 & <=1
    colWidth: float64 & <=1
    fiscalYearStartMonth: uint8 & <12 | *0
}

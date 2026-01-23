package constant_reference_discriminator

LayoutWithValue: GridLayoutUsingValue | RowsLayoutUsingValue
LayoutWithoutValue: GridLayoutWithoutValue | RowsLayoutWithoutValue

#GridLayoutKindType: "GridLayout"
#RowsLayoutKindType: "RowsLayout"

GridLayoutUsingValue: {
    kind: #GridLayoutKindType & "GridLayout"
    gridLayoutProperty: string
}

RowsLayoutUsingValue: {
    kind: #RowsLayoutKindType & "RowsLayout"
    rowsLayoutProperty: string
}

GridLayoutWithoutValue: {
    kind: #GridLayoutKindType
    gridLayoutProperty: string
}

RowsLayoutWithoutValue: {
    kind: #RowsLayoutKindType
    rowsLayoutProperty: string
}


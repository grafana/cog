package constant_reference_discriminator;


public class LayoutWithValue {
    protected GridLayoutUsingValue gridLayoutUsingValue;
    protected RowsLayoutUsingValue rowsLayoutUsingValue;
    protected LayoutWithValue() {}
    public static LayoutWithValue createGridLayoutUsingValue(GridLayoutUsingValue gridLayoutUsingValue) {
        LayoutWithValue layoutWithValue = new LayoutWithValue();
        layoutWithValue.gridLayoutUsingValue = gridLayoutUsingValue;
        return layoutWithValue;
    }
    public static LayoutWithValue createRowsLayoutUsingValue(RowsLayoutUsingValue rowsLayoutUsingValue) {
        LayoutWithValue layoutWithValue = new LayoutWithValue();
        layoutWithValue.rowsLayoutUsingValue = rowsLayoutUsingValue;
        return layoutWithValue;
    }
}

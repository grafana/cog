package constant_reference_discriminator;


public class LayoutWithoutValue {
    protected GridLayoutWithoutValue gridLayoutWithoutValue;
    protected RowsLayoutWithoutValue rowsLayoutWithoutValue;
    protected LayoutWithoutValue() {}
    public static LayoutWithoutValue createGridLayoutWithoutValue(GridLayoutWithoutValue gridLayoutWithoutValue) {
        LayoutWithoutValue layoutWithoutValue = new LayoutWithoutValue();
        layoutWithoutValue.gridLayoutWithoutValue = gridLayoutWithoutValue;
        return layoutWithoutValue;
    }
    public static LayoutWithoutValue createRowsLayoutWithoutValue(RowsLayoutWithoutValue rowsLayoutWithoutValue) {
        LayoutWithoutValue layoutWithoutValue = new LayoutWithoutValue();
        layoutWithoutValue.rowsLayoutWithoutValue = rowsLayoutWithoutValue;
        return layoutWithoutValue;
    }
}

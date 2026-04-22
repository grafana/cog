package constant_reference_discriminator;

import java.util.Objects;

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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof LayoutWithValue)) return false;
        LayoutWithValue o = (LayoutWithValue) other;
        if (!Objects.equals(this.gridLayoutUsingValue, o.gridLayoutUsingValue)) return false;
        if (!Objects.equals(this.rowsLayoutUsingValue, o.rowsLayoutUsingValue)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.gridLayoutUsingValue, this.rowsLayoutUsingValue);
    }
}

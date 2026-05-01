package constant_reference_discriminator;

import java.util.Objects;

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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof LayoutWithoutValue)) return false;
        LayoutWithoutValue o = (LayoutWithoutValue) other;
        if (!Objects.equals(this.gridLayoutWithoutValue, o.gridLayoutWithoutValue)) return false;
        if (!Objects.equals(this.rowsLayoutWithoutValue, o.rowsLayoutWithoutValue)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.gridLayoutWithoutValue, this.rowsLayoutWithoutValue);
    }
}

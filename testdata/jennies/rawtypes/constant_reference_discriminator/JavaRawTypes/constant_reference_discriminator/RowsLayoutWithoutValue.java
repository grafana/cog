package constant_reference_discriminator;

import java.util.Objects;

public class RowsLayoutWithoutValue {
    public String kind;
    public String rowsLayoutProperty;
    public RowsLayoutWithoutValue() {
        this.kind = Constants.RowsLayoutKindType;
        this.rowsLayoutProperty = "";
    }
    public RowsLayoutWithoutValue(String rowsLayoutProperty) {
        this.kind = Constants.RowsLayoutKindType;
        this.rowsLayoutProperty = rowsLayoutProperty;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof RowsLayoutWithoutValue)) return false;
        RowsLayoutWithoutValue o = (RowsLayoutWithoutValue) other;
        if (!Objects.equals(this.kind, o.kind)) return false;
        if (!Objects.equals(this.rowsLayoutProperty, o.rowsLayoutProperty)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.kind, this.rowsLayoutProperty);
    }
}

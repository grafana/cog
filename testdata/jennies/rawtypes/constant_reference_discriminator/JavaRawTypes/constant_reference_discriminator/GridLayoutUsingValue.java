package constant_reference_discriminator;

import java.util.Objects;

public class GridLayoutUsingValue {
    public String kind;
    public String gridLayoutProperty;
    public GridLayoutUsingValue() {
        this.kind = Constants.GridLayoutKindType;
        this.gridLayoutProperty = "";
    }
    public GridLayoutUsingValue(String gridLayoutProperty) {
        this.kind = Constants.GridLayoutKindType;
        this.gridLayoutProperty = gridLayoutProperty;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof GridLayoutUsingValue)) return false;
        GridLayoutUsingValue o = (GridLayoutUsingValue) other;
        if (!Objects.equals(this.kind, o.kind)) return false;
        if (!Objects.equals(this.gridLayoutProperty, o.gridLayoutProperty)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.kind, this.gridLayoutProperty);
    }
}

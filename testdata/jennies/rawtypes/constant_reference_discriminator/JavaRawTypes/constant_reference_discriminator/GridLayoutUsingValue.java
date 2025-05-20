package constant_reference_discriminator;


public class GridLayoutUsingValue {
    public String kind;
    public String gridLayoutProperty;
    public GridLayoutUsingValue() {
        this.kind = GridLayoutKindType;
        this.gridLayoutProperty = "";
    }
    public GridLayoutUsingValue(String gridLayoutProperty) {
        this.kind = GridLayoutKindType;
        this.gridLayoutProperty = gridLayoutProperty;
    }
}

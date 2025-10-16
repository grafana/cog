package constant_reference_discriminator;


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
}

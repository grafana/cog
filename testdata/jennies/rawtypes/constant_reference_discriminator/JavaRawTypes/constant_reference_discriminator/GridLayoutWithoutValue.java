package constant_reference_discriminator;


public class GridLayoutWithoutValue {
    public String kind;
    public String gridLayoutProperty;
    public GridLayoutWithoutValue() {
        this.kind = Constants.GridLayoutKindType;
        this.gridLayoutProperty = "";
    }
    public GridLayoutWithoutValue(String gridLayoutProperty) {
        this.kind = Constants.GridLayoutKindType;
        this.gridLayoutProperty = gridLayoutProperty;
    }
}

package constant_reference_discriminator;


public class GridLayoutWithoutValue {
    public String kind;
    public String gridLayoutProperty;
    public GridLayoutWithoutValue() {
        this.kind = GridLayoutKindType;
        this.gridLayoutProperty = "";
    }
    public GridLayoutWithoutValue(String gridLayoutProperty) {
        this.kind = GridLayoutKindType;
        this.gridLayoutProperty = gridLayoutProperty;
    }
}

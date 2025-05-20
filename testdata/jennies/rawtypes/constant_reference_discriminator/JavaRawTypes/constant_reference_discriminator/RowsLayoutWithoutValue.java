package constant_reference_discriminator;


public class RowsLayoutWithoutValue {
    public String kind;
    public String rowsLayoutProperty;
    public RowsLayoutWithoutValue() {
        this.kind = RowsLayoutKindType;
        this.rowsLayoutProperty = "";
    }
    public RowsLayoutWithoutValue(String rowsLayoutProperty) {
        this.kind = RowsLayoutKindType;
        this.rowsLayoutProperty = rowsLayoutProperty;
    }
}

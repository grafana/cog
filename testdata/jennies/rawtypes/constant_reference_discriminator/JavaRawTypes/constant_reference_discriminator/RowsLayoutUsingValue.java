package constant_reference_discriminator;


public class RowsLayoutUsingValue {
    public String kind;
    public String rowsLayoutProperty;
    public RowsLayoutUsingValue() {
        this.kind = RowsLayoutKindType;
    }
    public RowsLayoutUsingValue(String rowsLayoutProperty) {
        this.kind = RowsLayoutKindType;
        this.rowsLayoutProperty = rowsLayoutProperty;
    }
}

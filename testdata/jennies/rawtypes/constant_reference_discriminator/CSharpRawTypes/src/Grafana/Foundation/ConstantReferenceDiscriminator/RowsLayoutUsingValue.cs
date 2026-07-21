namespace Grafana.Foundation.ConstantReferenceDiscriminator;


public class RowsLayoutUsingValue
{
    public string Kind;
    public string RowsLayoutProperty;

    public RowsLayoutUsingValue()
    {
        this.Kind = default!;
        this.RowsLayoutProperty = "";
    }

    public RowsLayoutUsingValue(string kind, string rowsLayoutProperty)
    {
        this.Kind = kind;
        this.RowsLayoutProperty = rowsLayoutProperty;
    }
}

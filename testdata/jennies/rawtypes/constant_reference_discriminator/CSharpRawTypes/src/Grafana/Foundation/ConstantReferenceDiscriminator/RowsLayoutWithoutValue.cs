namespace Grafana.Foundation.ConstantReferenceDiscriminator;


public class RowsLayoutWithoutValue
{
    public string Kind;
    public string RowsLayoutProperty;

    public RowsLayoutWithoutValue()
    {
        this.Kind = default!;
        this.RowsLayoutProperty = "";
    }

    public RowsLayoutWithoutValue(string kind, string rowsLayoutProperty)
    {
        this.Kind = kind;
        this.RowsLayoutProperty = rowsLayoutProperty;
    }
}

namespace Grafana.Foundation.ConstantReferenceDiscriminator;


public class GridLayoutUsingValue
{
    public string Kind;
    public string GridLayoutProperty;

    public GridLayoutUsingValue()
    {
        this.Kind = default!;
        this.GridLayoutProperty = "";
    }

    public GridLayoutUsingValue(string kind, string gridLayoutProperty)
    {
        this.Kind = kind;
        this.GridLayoutProperty = gridLayoutProperty;
    }
}

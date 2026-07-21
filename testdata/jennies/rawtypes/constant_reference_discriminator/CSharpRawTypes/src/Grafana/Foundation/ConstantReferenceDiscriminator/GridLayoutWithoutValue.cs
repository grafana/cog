namespace Grafana.Foundation.ConstantReferenceDiscriminator;


public class GridLayoutWithoutValue
{
    public string Kind;
    public string GridLayoutProperty;

    public GridLayoutWithoutValue()
    {
        this.Kind = default!;
        this.GridLayoutProperty = "";
    }

    public GridLayoutWithoutValue(string kind, string gridLayoutProperty)
    {
        this.Kind = kind;
        this.GridLayoutProperty = gridLayoutProperty;
    }
}

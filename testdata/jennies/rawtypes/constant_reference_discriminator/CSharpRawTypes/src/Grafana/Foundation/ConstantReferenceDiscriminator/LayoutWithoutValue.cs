namespace Grafana.Foundation.ConstantReferenceDiscriminator;


public class LayoutWithoutValue
{
    public GridLayoutWithoutValue GridLayoutWithoutValue;
    public RowsLayoutWithoutValue RowsLayoutWithoutValue;

    public LayoutWithoutValue()
    {
    }

    public LayoutWithoutValue(GridLayoutWithoutValue gridLayoutWithoutValue, RowsLayoutWithoutValue rowsLayoutWithoutValue)
    {
        this.GridLayoutWithoutValue = gridLayoutWithoutValue;
        this.RowsLayoutWithoutValue = rowsLayoutWithoutValue;
    }
}

namespace Grafana.Foundation.ConstantReferenceDiscriminator;


public class LayoutWithValue
{
    public GridLayoutUsingValue GridLayoutUsingValue;
    public RowsLayoutUsingValue RowsLayoutUsingValue;

    public LayoutWithValue()
    {
    }

    public LayoutWithValue(GridLayoutUsingValue gridLayoutUsingValue, RowsLayoutUsingValue rowsLayoutUsingValue)
    {
        this.GridLayoutUsingValue = gridLayoutUsingValue;
        this.RowsLayoutUsingValue = rowsLayoutUsingValue;
    }
}

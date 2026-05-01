namespace Grafana.Foundation.Disjunctions;


public class SeveralRefs
{
    public SomeStruct SomeStruct;
    public SomeOtherStruct SomeOtherStruct;
    public YetAnotherStruct YetAnotherStruct;

    public SeveralRefs()
    {
    }

    public SeveralRefs(SomeStruct someStruct, SomeOtherStruct someOtherStruct, YetAnotherStruct yetAnotherStruct)
    {
        this.SomeStruct = someStruct;
        this.SomeOtherStruct = someOtherStruct;
        this.YetAnotherStruct = yetAnotherStruct;
    }
}

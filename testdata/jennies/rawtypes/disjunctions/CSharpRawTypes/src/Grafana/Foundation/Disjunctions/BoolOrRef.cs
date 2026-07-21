namespace Grafana.Foundation.Disjunctions;


public class BoolOrRef
{
    public bool Bool;
    public SomeStruct SomeStruct;

    public BoolOrRef()
    {
    }

    public BoolOrRef(bool boolArg, SomeStruct someStruct)
    {
        this.Bool = boolArg;
        this.SomeStruct = someStruct;
    }
}

namespace Grafana.Foundation.StructComplexFields;


public class StringOrSomeOtherStruct
{
    public string String;
    public SomeOtherStruct SomeOtherStruct;

    public StringOrSomeOtherStruct()
    {
    }

    public StringOrSomeOtherStruct(string stringArg, SomeOtherStruct someOtherStruct)
    {
        this.String = stringArg;
        this.SomeOtherStruct = someOtherStruct;
    }
}

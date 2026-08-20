namespace Grafana.Foundation.StructComplexFields;


public class StringOrBool
{
    public string String;
    public bool Bool;

    public StringOrBool()
    {
    }

    public StringOrBool(string stringArg, bool boolArg)
    {
        this.String = stringArg;
        this.Bool = boolArg;
    }
}

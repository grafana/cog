namespace Grafana.Foundation.Defaults;

using System.Collections.Generic;

public class DefaultsStructComplexField
{
    public string Uid;
    public DefaultsStructComplexFieldNested Nested;
    public List<string> Array;

    public DefaultsStructComplexField()
    {
        this.Uid = "";
        this.Nested = new DefaultsStructComplexFieldNested();
        this.Array = new List<string>();
    }

    public DefaultsStructComplexField(string uid, DefaultsStructComplexFieldNested nested, List<string> array)
    {
        this.Uid = uid;
        this.Nested = nested;
        this.Array = array;
    }
}

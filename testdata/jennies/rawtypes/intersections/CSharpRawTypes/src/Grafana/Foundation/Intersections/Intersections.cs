namespace Grafana.Foundation.Intersections;

using Grafana.Foundation.ExternalPkg;

public class Intersections : SomeStruct, AnotherStruct
{
    public string FieldString;
    public int FieldInteger;

    public Intersections()
    {
        this.FieldString = "hello";
        this.FieldInteger = 32;
    }

    public Intersections(string fieldString, int fieldInteger)
    {
        this.FieldString = fieldString;
        this.FieldInteger = fieldInteger;
    }
}

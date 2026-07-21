namespace Grafana.Foundation.Disjunctions;


public class SomeStruct
{
    public string Type;
    public object FieldAny;

    public SomeStruct()
    {
        this.Type = "";
        this.FieldAny = new object();
    }

    public SomeStruct(string type, object fieldAny)
    {
        this.Type = type;
        this.FieldAny = fieldAny;
    }
}

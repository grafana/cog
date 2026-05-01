namespace Grafana.Foundation.Maps;


public class SomeStruct
{
    public object FieldAny;

    public SomeStruct()
    {
        this.FieldAny = new object();
    }

    public SomeStruct(object fieldAny)
    {
        this.FieldAny = fieldAny;
    }
}

namespace Grafana.Foundation.Refs;


public class RefToSomeStruct
{
    public object FieldAny;

    public RefToSomeStruct()
    {
        this.FieldAny = new object();
    }

    public RefToSomeStruct(object fieldAny)
    {
        this.FieldAny = fieldAny;
    }
}

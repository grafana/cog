namespace Grafana.Foundation.StructComplexFields;


public class SomeOtherStruct
{
    public object FieldAny;

    public SomeOtherStruct()
    {
        this.FieldAny = new object();
    }

    public SomeOtherStruct(object fieldAny)
    {
        this.FieldAny = fieldAny;
    }
}

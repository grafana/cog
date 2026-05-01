namespace Grafana.Foundation.Intersections;


public class SomeStruct
{
    public bool FieldBool;

    public SomeStruct()
    {
        this.FieldBool = true;
    }

    public SomeStruct(bool fieldBool)
    {
        this.FieldBool = fieldBool;
    }
}

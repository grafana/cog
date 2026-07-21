namespace Grafana.Foundation.NullableFields;


public class MyObject
{
    public string Field;

    public MyObject()
    {
        this.Field = "";
    }

    public MyObject(string field)
    {
        this.Field = field;
    }
}

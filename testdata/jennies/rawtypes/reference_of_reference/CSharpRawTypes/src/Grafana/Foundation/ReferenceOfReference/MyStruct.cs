namespace Grafana.Foundation.ReferenceOfReference;


public class MyStruct
{
    public OtherStruct Field;

    public MyStruct()
    {
    }

    public MyStruct(OtherStruct field)
    {
        this.Field = field;
    }
}

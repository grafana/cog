namespace Grafana.Foundation.Disjunctions;


public class SomeOtherStruct
{
    public string Type;
    public byte Foo;

    public SomeOtherStruct()
    {
        this.Type = "";
        this.Foo = (byte) 0;
    }

    public SomeOtherStruct(string type, byte foo)
    {
        this.Type = type;
        this.Foo = foo;
    }
}

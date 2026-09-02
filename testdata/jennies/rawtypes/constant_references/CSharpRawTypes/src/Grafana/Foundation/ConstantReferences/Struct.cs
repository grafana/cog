namespace Grafana.Foundation.ConstantReferences;


public class Struct
{
    public string MyValue;
    public Enum MyEnum;

    public Struct()
    {
        this.MyValue = "";
        this.MyEnum = Enum.ValueA;
    }

    public Struct(string myValue, Enum myEnum)
    {
        this.MyValue = myValue;
        this.MyEnum = myEnum;
    }
}

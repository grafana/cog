namespace Grafana.Foundation.ConstantReferences;


public class StructB
{
    public Enum MyEnum;
    public string MyValue;

    public StructB()
    {
        this.MyEnum = default!;
        this.MyValue = "";
    }

    public StructB(Enum myEnum, string myValue)
    {
        this.MyEnum = myEnum;
        this.MyValue = myValue;
    }
}

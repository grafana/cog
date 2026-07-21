namespace Grafana.Foundation.ConstantReferences;


public class StructA
{
    public Enum MyEnum;
    public Enum Other;

    public StructA()
    {
        this.MyEnum = default!;
    }

    public StructA(Enum myEnum, Enum other)
    {
        this.MyEnum = myEnum;
        this.Other = other;
    }
}

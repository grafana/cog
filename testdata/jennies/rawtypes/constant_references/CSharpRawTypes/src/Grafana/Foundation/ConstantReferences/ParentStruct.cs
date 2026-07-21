namespace Grafana.Foundation.ConstantReferences;


public class ParentStruct
{
    public Enum MyEnum;

    public ParentStruct()
    {
        this.MyEnum = Enum.ValueA;
    }

    public ParentStruct(Enum myEnum)
    {
        this.MyEnum = myEnum;
    }
}

namespace Grafana.Foundation.Defaults;


public class SomeStruct
{
    public bool FieldBool;
    public string FieldString;
    public string FieldStringWithConstantValue;
    public float FieldFloat32;
    public int FieldInt32;

    public SomeStruct()
    {
        this.FieldBool = true;
        this.FieldString = "foo";
        this.FieldStringWithConstantValue = "";
        this.FieldFloat32 = 42.42f;
        this.FieldInt32 = 42;
    }

    public SomeStruct(bool fieldBool, string fieldString, string fieldStringWithConstantValue, float fieldFloat32, int fieldInt32)
    {
        this.FieldBool = fieldBool;
        this.FieldString = fieldString;
        this.FieldStringWithConstantValue = fieldStringWithConstantValue;
        this.FieldFloat32 = fieldFloat32;
        this.FieldInt32 = fieldInt32;
    }
}

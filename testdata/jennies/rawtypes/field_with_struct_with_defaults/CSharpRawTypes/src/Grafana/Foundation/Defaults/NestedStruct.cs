namespace Grafana.Foundation.Defaults;


public class NestedStruct
{
    public string StringVal;
    public long IntVal;

    public NestedStruct()
    {
        this.StringVal = "";
        this.IntVal = 0L;
    }

    public NestedStruct(string stringVal, long intVal)
    {
        this.StringVal = stringVal;
        this.IntVal = intVal;
    }
}

namespace Grafana.Foundation.Defaults;


public class DefaultsStructPartialComplexField
{
    public string Uid;
    public long IntVal;

    public DefaultsStructPartialComplexField()
    {
        this.Uid = "";
        this.IntVal = 0L;
    }

    public DefaultsStructPartialComplexField(string uid, long intVal)
    {
        this.Uid = uid;
        this.IntVal = intVal;
    }
}

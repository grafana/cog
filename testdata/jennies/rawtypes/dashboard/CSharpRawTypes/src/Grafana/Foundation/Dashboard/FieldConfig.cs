namespace Grafana.Foundation.Dashboard;


public class FieldConfig
{
    public string Unit;
    public object Custom;

    public FieldConfig()
    {
    }

    public FieldConfig(string unit, object custom)
    {
        this.Unit = unit;
        this.Custom = custom;
    }
}

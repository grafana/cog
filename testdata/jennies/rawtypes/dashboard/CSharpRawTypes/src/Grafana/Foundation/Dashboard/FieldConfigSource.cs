namespace Grafana.Foundation.Dashboard;


public class FieldConfigSource
{
    public FieldConfig Defaults;

    public FieldConfigSource()
    {
    }

    public FieldConfigSource(FieldConfig defaults)
    {
        this.Defaults = defaults;
    }
}

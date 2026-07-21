namespace Grafana.Foundation.VariantPanelcfgFull;


public class FieldConfig
{
    public string TimeseriesFieldConfigOption;

    public FieldConfig()
    {
        this.TimeseriesFieldConfigOption = "";
    }

    public FieldConfig(string timeseriesFieldConfigOption)
    {
        this.TimeseriesFieldConfigOption = timeseriesFieldConfigOption;
    }
}

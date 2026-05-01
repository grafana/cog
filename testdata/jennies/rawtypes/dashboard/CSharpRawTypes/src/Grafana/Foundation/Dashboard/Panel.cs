namespace Grafana.Foundation.Dashboard;

using System.Collections.Generic;

public class Panel
{
    public string Title;
    public string Type;
    public DataSourceRef Datasource;
    public object Options;
    public List<Cog.Variants.Dataquery> Targets;
    public FieldConfigSource FieldConfig;

    public Panel()
    {
        this.Title = "";
        this.Type = "";
    }

    public Panel(string title, string type, DataSourceRef datasource, object options, List<Cog.Variants.Dataquery> targets, FieldConfigSource fieldConfig)
    {
        this.Title = title;
        this.Type = type;
        this.Datasource = datasource;
        this.Options = options;
        this.Targets = targets;
        this.FieldConfig = fieldConfig;
    }
}

namespace Grafana.Foundation.Dashboard;


public class DataSourceRef
{
    public string Type;
    public string Uid;

    public DataSourceRef()
    {
    }

    public DataSourceRef(string type, string uid)
    {
        this.Type = type;
        this.Uid = uid;
    }
}

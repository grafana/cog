namespace Grafana.Foundation.Intersections;


// Base properties for all metrics
public class Common
{
    // The metric name
    public string Name;
    // The metric type
    public CommonType Type;
    // The type of data the metric contains
    public CommonContains Contains;

    public Common()
    {
        this.Name = "";
        this.Type = CommonType.Counter;
        this.Contains = CommonContains.Default;
    }

    public Common(string name, CommonType type, CommonContains contains)
    {
        this.Name = name;
        this.Type = type;
        this.Contains = contains;
    }
}

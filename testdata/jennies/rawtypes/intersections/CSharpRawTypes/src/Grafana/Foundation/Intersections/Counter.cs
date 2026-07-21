namespace Grafana.Foundation.Intersections;


// Counter metric combining common properties with specific values
public class Counter : Common
{
    // Counter metric values
    public object Values;

    public Counter()
    {
        this.Values = new object();
    }

    public Counter(object values)
    {
        this.Values = values;
    }
}

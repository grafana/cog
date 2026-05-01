namespace Grafana.Foundation.TimeHint;


public class ObjWithTimeField
{
    public string RegisteredAt;
    public string Duration;

    public ObjWithTimeField()
    {
        this.RegisteredAt = "";
        this.Duration = "";
    }

    public ObjWithTimeField(string registeredAt, string duration)
    {
        this.RegisteredAt = registeredAt;
        this.Duration = duration;
    }
}

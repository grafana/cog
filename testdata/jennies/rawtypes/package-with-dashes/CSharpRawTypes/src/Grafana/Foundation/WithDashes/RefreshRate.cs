namespace Grafana.Foundation.WithDashes;


// Refresh rate or disabled.
public class RefreshRate
{
    public string String;
    public bool Bool;

    public RefreshRate()
    {
    }

    public RefreshRate(string stringArg, bool boolArg)
    {
        this.String = stringArg;
        this.Bool = boolArg;
    }
}

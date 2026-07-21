namespace Grafana.Foundation.Disjunctions;


public class YetAnotherStruct
{
    public string Type;
    public byte Bar;

    public YetAnotherStruct()
    {
        this.Type = "";
        this.Bar = 0;
    }

    public YetAnotherStruct(string type, byte bar)
    {
        this.Type = type;
        this.Bar = bar;
    }
}

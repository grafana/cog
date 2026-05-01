namespace Grafana.Foundation.NullableFields;

using System.Collections.Generic;

public class Struct
{
    public MyObject A;
    public MyObject B;
    public string C;
    public List<string> D;
    public Dictionary<string, string> E;
    public NullableFieldsStructF F;
    public string G;

    public Struct()
    {
        this.E = new Dictionary<string, string>();
    }

    public Struct(MyObject a, MyObject b, string c, List<string> d, Dictionary<string, string> e, NullableFieldsStructF f, string g)
    {
        this.A = a;
        this.B = b;
        this.C = c;
        this.D = d;
        this.E = e;
        this.F = f;
        this.G = g;
    }
}

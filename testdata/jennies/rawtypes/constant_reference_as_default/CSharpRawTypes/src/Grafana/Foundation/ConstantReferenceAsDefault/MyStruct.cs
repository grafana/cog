namespace Grafana.Foundation.ConstantReferenceAsDefault;


public class MyStruct
{
    public string AString;
    public string OptString;

    public MyStruct()
    {
        this.AString = default!;
    }

    public MyStruct(string aString, string optString)
    {
        this.AString = aString;
        this.OptString = optString;
    }
}

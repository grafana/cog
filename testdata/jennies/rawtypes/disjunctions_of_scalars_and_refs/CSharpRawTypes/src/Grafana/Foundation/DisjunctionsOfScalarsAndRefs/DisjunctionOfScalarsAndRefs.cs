namespace Grafana.Foundation.DisjunctionsOfScalarsAndRefs;

using System.Collections.Generic;

public class DisjunctionOfScalarsAndRefs
{
    public string String;
    public bool Bool;
    public List<string> ArrayOfString;
    public MyRefA MyRefA;
    public MyRefB MyRefB;

    public DisjunctionOfScalarsAndRefs()
    {
    }

    public DisjunctionOfScalarsAndRefs(string stringArg, bool boolArg, List<string> arrayOfString, MyRefA myRefA, MyRefB myRefB)
    {
        this.String = stringArg;
        this.Bool = boolArg;
        this.ArrayOfString = arrayOfString;
        this.MyRefA = myRefA;
        this.MyRefB = myRefB;
    }
}

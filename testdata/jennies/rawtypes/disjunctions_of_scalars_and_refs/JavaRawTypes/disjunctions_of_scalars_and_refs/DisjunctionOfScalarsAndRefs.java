package disjunctions_of_scalars_and_refs;

import java.util.Objects;
import java.util.List;

public class DisjunctionOfScalarsAndRefs {
    protected String string;
    protected Boolean bool;
    protected List<String> arrayOfString;
    protected MyRefA myRefA;
    protected MyRefB myRefB;
    protected DisjunctionOfScalarsAndRefs() {}
    public static DisjunctionOfScalarsAndRefs createString(String string) {
        DisjunctionOfScalarsAndRefs disjunctionOfScalarsAndRefs = new DisjunctionOfScalarsAndRefs();
        disjunctionOfScalarsAndRefs.string = string;
        return disjunctionOfScalarsAndRefs;
    }
    public static DisjunctionOfScalarsAndRefs createBool(Boolean bool) {
        DisjunctionOfScalarsAndRefs disjunctionOfScalarsAndRefs = new DisjunctionOfScalarsAndRefs();
        disjunctionOfScalarsAndRefs.bool = bool;
        return disjunctionOfScalarsAndRefs;
    }
    public static DisjunctionOfScalarsAndRefs createArrayOfString(List<String> arrayOfString) {
        DisjunctionOfScalarsAndRefs disjunctionOfScalarsAndRefs = new DisjunctionOfScalarsAndRefs();
        disjunctionOfScalarsAndRefs.arrayOfString = arrayOfString;
        return disjunctionOfScalarsAndRefs;
    }
    public static DisjunctionOfScalarsAndRefs createMyRefA(MyRefA myRefA) {
        DisjunctionOfScalarsAndRefs disjunctionOfScalarsAndRefs = new DisjunctionOfScalarsAndRefs();
        disjunctionOfScalarsAndRefs.myRefA = myRefA;
        return disjunctionOfScalarsAndRefs;
    }
    public static DisjunctionOfScalarsAndRefs createMyRefB(MyRefB myRefB) {
        DisjunctionOfScalarsAndRefs disjunctionOfScalarsAndRefs = new DisjunctionOfScalarsAndRefs();
        disjunctionOfScalarsAndRefs.myRefB = myRefB;
        return disjunctionOfScalarsAndRefs;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof DisjunctionOfScalarsAndRefs)) return false;
        DisjunctionOfScalarsAndRefs o = (DisjunctionOfScalarsAndRefs) other;
        if (!Objects.equals(this.string, o.string)) return false;
        if (!Objects.equals(this.bool, o.bool)) return false;
        if (!Objects.equals(this.arrayOfString, o.arrayOfString)) return false;
        if (!Objects.equals(this.myRefA, o.myRefA)) return false;
        if (!Objects.equals(this.myRefB, o.myRefB)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.string, this.bool, this.arrayOfString, this.myRefA, this.myRefB);
    }
}

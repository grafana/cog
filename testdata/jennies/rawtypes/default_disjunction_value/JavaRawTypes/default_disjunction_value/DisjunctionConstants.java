package default_disjunction_value;


public class DisjunctionConstants {
    protected String string;
    protected Long int64;
    protected Boolean bool;
    protected DisjunctionConstants() {}
    public static DisjunctionConstants createString(String string) {
        DisjunctionConstants disjunctionConstants = new DisjunctionConstants();
        disjunctionConstants.string = string;
        return disjunctionConstants;
    }
    public static DisjunctionConstants createInt64(Long int64) {
        DisjunctionConstants disjunctionConstants = new DisjunctionConstants();
        disjunctionConstants.int64 = int64;
        return disjunctionConstants;
    }
    public static DisjunctionConstants createBool(Boolean bool) {
        DisjunctionConstants disjunctionConstants = new DisjunctionConstants();
        disjunctionConstants.bool = bool;
        return disjunctionConstants;
    }
}

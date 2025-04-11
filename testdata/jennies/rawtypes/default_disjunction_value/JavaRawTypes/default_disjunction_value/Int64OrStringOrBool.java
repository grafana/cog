package default_disjunction_value;


public class Int64OrStringOrBool {
    protected Long int64;
    protected String string;
    protected Boolean bool;
    protected Int64OrStringOrBool() {}
    public static Int64OrStringOrBool createInt64(Long int64) {
        Int64OrStringOrBool int64OrStringOrBool = new Int64OrStringOrBool();
        int64OrStringOrBool.int64 = int64;
        return int64OrStringOrBool;
    }
    public static Int64OrStringOrBool createString(String string) {
        Int64OrStringOrBool int64OrStringOrBool = new Int64OrStringOrBool();
        int64OrStringOrBool.string = string;
        return int64OrStringOrBool;
    }
    public static Int64OrStringOrBool createBool(Boolean bool) {
        Int64OrStringOrBool int64OrStringOrBool = new Int64OrStringOrBool();
        int64OrStringOrBool.bool = bool;
        return int64OrStringOrBool;
    }
}

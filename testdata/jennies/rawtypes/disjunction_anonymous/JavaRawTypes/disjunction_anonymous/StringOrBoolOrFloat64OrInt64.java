package disjunction_anonymous;


public class StringOrBoolOrFloat64OrInt64 {
    protected String string;
    protected Boolean bool;
    protected Double float64;
    protected Long int64;
    protected StringOrBoolOrFloat64OrInt64() {}
    public static StringOrBoolOrFloat64OrInt64 createString(String string) {
        StringOrBoolOrFloat64OrInt64 stringOrBoolOrFloat64OrInt64 = new StringOrBoolOrFloat64OrInt64();
        stringOrBoolOrFloat64OrInt64.string = string;
        return stringOrBoolOrFloat64OrInt64;
    }
    public static StringOrBoolOrFloat64OrInt64 createBool(Boolean bool) {
        StringOrBoolOrFloat64OrInt64 stringOrBoolOrFloat64OrInt64 = new StringOrBoolOrFloat64OrInt64();
        stringOrBoolOrFloat64OrInt64.bool = bool;
        return stringOrBoolOrFloat64OrInt64;
    }
    public static StringOrBoolOrFloat64OrInt64 createFloat64(Double float64) {
        StringOrBoolOrFloat64OrInt64 stringOrBoolOrFloat64OrInt64 = new StringOrBoolOrFloat64OrInt64();
        stringOrBoolOrFloat64OrInt64.float64 = float64;
        return stringOrBoolOrFloat64OrInt64;
    }
    public static StringOrBoolOrFloat64OrInt64 createInt64(Long int64) {
        StringOrBoolOrFloat64OrInt64 stringOrBoolOrFloat64OrInt64 = new StringOrBoolOrFloat64OrInt64();
        stringOrBoolOrFloat64OrInt64.int64 = int64;
        return stringOrBoolOrFloat64OrInt64;
    }
}

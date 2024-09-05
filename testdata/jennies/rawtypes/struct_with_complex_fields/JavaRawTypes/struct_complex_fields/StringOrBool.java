package struct_complex_fields;


public class StringOrBool {
    protected String string;
    protected Boolean bool;
    protected StringOrBool() {}
    public static StringOrBool createString(String string) {
        StringOrBool stringOrBool = new StringOrBool();
        stringOrBool.string = string;
        return stringOrBool;
    }
    public static StringOrBool createBool(Boolean bool) {
        StringOrBool stringOrBool = new StringOrBool();
        stringOrBool.bool = bool;
        return stringOrBool;
    }
}

package struct_complex_fields;


public class StringOrSomeOtherStruct {
    protected String string;
    protected SomeOtherStruct someOtherStruct;
    protected StringOrSomeOtherStruct() {}
    public static StringOrSomeOtherStruct createString(String string) {
        StringOrSomeOtherStruct stringOrSomeOtherStruct = new StringOrSomeOtherStruct();
        stringOrSomeOtherStruct.string = string;
        return stringOrSomeOtherStruct;
    }
    public static StringOrSomeOtherStruct createSomeOtherStruct(SomeOtherStruct someOtherStruct) {
        StringOrSomeOtherStruct stringOrSomeOtherStruct = new StringOrSomeOtherStruct();
        stringOrSomeOtherStruct.someOtherStruct = someOtherStruct;
        return stringOrSomeOtherStruct;
    }
}

package struct_complex_fields;


public class StringOrSomeOtherStruct {
    public String string;
    public SomeOtherStruct someOtherStruct;
    public StringOrSomeOtherStruct() {
    }
    public StringOrSomeOtherStruct(String string,SomeOtherStruct someOtherStruct) {
        this.string = string;
        this.someOtherStruct = someOtherStruct;
    }
}

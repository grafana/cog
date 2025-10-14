package constant_reference_as_default;


public class MyStruct {
    public String aString;
    public String optString;
    public MyStruct() {
        this.aString = Constants.ConstantRefString;
        this.optString = Constants.ConstantRefString;
    }
}

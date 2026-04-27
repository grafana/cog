package reference_of_reference;


public class MyStruct {
    public OtherStruct field;
    public MyStruct() {
    }
    public MyStruct(OtherStruct field) {
        this.field = field;
    }
}

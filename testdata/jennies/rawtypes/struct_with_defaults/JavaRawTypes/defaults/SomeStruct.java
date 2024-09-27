package defaults;


public class SomeStruct {
    public Boolean fieldBool;
    public String fieldString;
    public String fieldStringWithConstantValue;
    public Float fieldFloat32;
    public Integer fieldInt32;

    public SomeStruct() {
        this.fieldBool = true;
        this.fieldString = "foo";
        this.fieldFloat32 = 42.4f;
        this.fieldInt32 = 42;
    }
}

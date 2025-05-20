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
        this.fieldStringWithConstantValue = "";
        this.fieldFloat32 = 42.4f;
        this.fieldInt32 = 42;
    }
    public SomeStruct(Boolean fieldBool,String fieldString,String fieldStringWithConstantValue,Float fieldFloat32,Integer fieldInt32) {
        this.fieldBool = fieldBool;
        this.fieldString = fieldString;
        this.fieldStringWithConstantValue = fieldStringWithConstantValue;
        this.fieldFloat32 = fieldFloat32;
        this.fieldInt32 = fieldInt32;
    }
}

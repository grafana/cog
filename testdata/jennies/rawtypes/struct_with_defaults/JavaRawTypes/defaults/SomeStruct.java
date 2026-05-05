package defaults;

import java.util.Objects;

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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldBool, o.fieldBool)) return false;
        if (!Objects.equals(this.fieldString, o.fieldString)) return false;
        if (!Objects.equals(this.fieldStringWithConstantValue, o.fieldStringWithConstantValue)) return false;
        if (!Objects.equals(this.fieldFloat32, o.fieldFloat32)) return false;
        if (!Objects.equals(this.fieldInt32, o.fieldInt32)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldBool, this.fieldString, this.fieldStringWithConstantValue, this.fieldFloat32, this.fieldInt32);
    }
}

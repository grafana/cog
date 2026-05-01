package intersections;

import java.util.Objects;

public class SomeStruct {
    public Boolean fieldBool;
    public SomeStruct() {
        this.fieldBool = true;
    }
    public SomeStruct(Boolean fieldBool) {
        this.fieldBool = fieldBool;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldBool, o.fieldBool)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldBool);
    }
}

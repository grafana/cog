package struct_complex_fields;

import java.util.Objects;

public class SomeOtherStruct {
    public Object fieldAny;
    public SomeOtherStruct() {
        this.fieldAny = new Object();
    }
    public SomeOtherStruct(Object fieldAny) {
        this.fieldAny = fieldAny;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeOtherStruct)) return false;
        SomeOtherStruct o = (SomeOtherStruct) other;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldAny);
    }
}

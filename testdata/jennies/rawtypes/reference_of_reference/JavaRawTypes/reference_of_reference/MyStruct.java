package reference_of_reference;

import java.util.Objects;

public class MyStruct {
    public OtherStruct field;
    public MyStruct() {
    }
    public MyStruct(OtherStruct field) {
        this.field = field;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof MyStruct)) return false;
        MyStruct o = (MyStruct) other;
        if (!Objects.equals(this.field, o.field)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.field);
    }
}

package disjunctions;

import java.util.Objects;

public class SomeOtherStruct {
    public String type;
    public Byte foo;
    public SomeOtherStruct() {
        this.type = "";
        this.foo = (byte) 0;
    }
    public SomeOtherStruct(String type,Byte foo) {
        this.type = type;
        this.foo = foo;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeOtherStruct)) return false;
        SomeOtherStruct o = (SomeOtherStruct) other;
        if (!Objects.equals(this.type, o.type)) return false;
        if (!Objects.equals(this.foo, o.foo)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.type, this.foo);
    }
}

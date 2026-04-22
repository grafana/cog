package disjunctions;

import java.util.Objects;

public class BoolOrRef {
    protected Boolean bool;
    protected SomeStruct someStruct;
    protected BoolOrRef() {}
    public static BoolOrRef createBool(Boolean bool) {
        BoolOrRef boolOrRef = new BoolOrRef();
        boolOrRef.bool = bool;
        return boolOrRef;
    }
    public static BoolOrRef createSomeStruct(SomeStruct someStruct) {
        BoolOrRef boolOrRef = new BoolOrRef();
        boolOrRef.someStruct = someStruct;
        return boolOrRef;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof BoolOrRef)) return false;
        BoolOrRef o = (BoolOrRef) other;
        if (!Objects.equals(this.bool, o.bool)) return false;
        if (!Objects.equals(this.someStruct, o.someStruct)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.bool, this.someStruct);
    }
}

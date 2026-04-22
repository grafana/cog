package disjunctions;

import java.util.Objects;

public class SeveralRefs {
    protected SomeStruct someStruct;
    protected SomeOtherStruct someOtherStruct;
    protected YetAnotherStruct yetAnotherStruct;
    protected SeveralRefs() {}
    public static SeveralRefs createSomeStruct(SomeStruct someStruct) {
        SeveralRefs severalRefs = new SeveralRefs();
        severalRefs.someStruct = someStruct;
        return severalRefs;
    }
    public static SeveralRefs createSomeOtherStruct(SomeOtherStruct someOtherStruct) {
        SeveralRefs severalRefs = new SeveralRefs();
        severalRefs.someOtherStruct = someOtherStruct;
        return severalRefs;
    }
    public static SeveralRefs createYetAnotherStruct(YetAnotherStruct yetAnotherStruct) {
        SeveralRefs severalRefs = new SeveralRefs();
        severalRefs.yetAnotherStruct = yetAnotherStruct;
        return severalRefs;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SeveralRefs)) return false;
        SeveralRefs o = (SeveralRefs) other;
        if (!Objects.equals(this.someStruct, o.someStruct)) return false;
        if (!Objects.equals(this.someOtherStruct, o.someOtherStruct)) return false;
        if (!Objects.equals(this.yetAnotherStruct, o.yetAnotherStruct)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.someStruct, this.someOtherStruct, this.yetAnotherStruct);
    }
}

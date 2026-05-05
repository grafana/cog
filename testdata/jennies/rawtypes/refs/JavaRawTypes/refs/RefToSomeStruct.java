package refs;

import java.util.Objects;

public class RefToSomeStruct {
    public Object fieldAny;
    public RefToSomeStruct() {
        this.fieldAny = new Object();
    }
    public RefToSomeStruct(Object fieldAny) {
        this.fieldAny = fieldAny;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof RefToSomeStruct)) return false;
        RefToSomeStruct o = (RefToSomeStruct) other;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldAny);
    }
}

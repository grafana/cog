package maps;

import java.util.Objects;

public class SomeStruct {
    public Object fieldAny;
    public SomeStruct() {
        this.fieldAny = new Object();
    }
    public SomeStruct(Object fieldAny) {
        this.fieldAny = fieldAny;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldAny);
    }
}

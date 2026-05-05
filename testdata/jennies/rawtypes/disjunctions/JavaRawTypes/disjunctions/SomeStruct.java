package disjunctions;

import java.util.Objects;

public class SomeStruct {
    public String type;
    public Object fieldAny;
    public SomeStruct() {
        this.type = "";
        this.fieldAny = new Object();
    }
    public SomeStruct(String type,Object fieldAny) {
        this.type = type;
        this.fieldAny = fieldAny;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.type, o.type)) return false;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.type, this.fieldAny);
    }
}

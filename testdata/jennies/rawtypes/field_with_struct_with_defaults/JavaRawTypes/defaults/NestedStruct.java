package defaults;

import java.util.Objects;

public class NestedStruct {
    public String stringVal;
    public Long intVal;
    public NestedStruct() {
        this.stringVal = "";
        this.intVal = 0L;
    }
    public NestedStruct(String stringVal,Long intVal) {
        this.stringVal = stringVal;
        this.intVal = intVal;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof NestedStruct)) return false;
        NestedStruct o = (NestedStruct) other;
        if (!Objects.equals(this.stringVal, o.stringVal)) return false;
        if (!Objects.equals(this.intVal, o.intVal)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.stringVal, this.intVal);
    }
}

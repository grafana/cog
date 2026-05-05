package defaults;

import java.util.Objects;

public class DefaultsStructPartialComplexField {
    public String uid;
    public Long intVal;
    public DefaultsStructPartialComplexField() {
        this.uid = "";
        this.intVal = 0L;
    }
    public DefaultsStructPartialComplexField(String uid,Long intVal) {
        this.uid = uid;
        this.intVal = intVal;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof DefaultsStructPartialComplexField)) return false;
        DefaultsStructPartialComplexField o = (DefaultsStructPartialComplexField) other;
        if (!Objects.equals(this.uid, o.uid)) return false;
        if (!Objects.equals(this.intVal, o.intVal)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.uid, this.intVal);
    }
}

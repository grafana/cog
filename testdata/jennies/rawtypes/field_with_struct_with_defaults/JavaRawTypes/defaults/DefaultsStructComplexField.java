package defaults;

import java.util.Objects;
import java.util.LinkedList;
import java.util.List;

public class DefaultsStructComplexField {
    public String uid;
    public DefaultsStructComplexFieldNested nested;
    public List<String> array;
    public DefaultsStructComplexField() {
        this.uid = "";
        this.nested = new defaults.DefaultsStructComplexFieldNested();
        this.array = new LinkedList<>();
    }
    public DefaultsStructComplexField(String uid,DefaultsStructComplexFieldNested nested,List<String> array) {
        this.uid = uid;
        this.nested = nested;
        this.array = array;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof DefaultsStructComplexField)) return false;
        DefaultsStructComplexField o = (DefaultsStructComplexField) other;
        if (!Objects.equals(this.uid, o.uid)) return false;
        if (!Objects.equals(this.nested, o.nested)) return false;
        if (!Objects.equals(this.array, o.array)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.uid, this.nested, this.array);
    }
}

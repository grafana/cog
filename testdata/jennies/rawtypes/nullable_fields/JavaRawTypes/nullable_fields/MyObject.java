package nullable_fields;

import java.util.Objects;

public class MyObject {
    public String field;
    public MyObject() {
        this.field = "";
    }
    public MyObject(String field) {
        this.field = field;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof MyObject)) return false;
        MyObject o = (MyObject) other;
        if (!Objects.equals(this.field, o.field)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.field);
    }
}

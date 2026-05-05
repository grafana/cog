package nullable_fields;

import java.util.Objects;

public class NullableFieldsStructF {
    public String a;
    public NullableFieldsStructF() {
        this.a = "";
    }
    public NullableFieldsStructF(String a) {
        this.a = a;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof NullableFieldsStructF)) return false;
        NullableFieldsStructF o = (NullableFieldsStructF) other;
        if (!Objects.equals(this.a, o.a)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.a);
    }
}

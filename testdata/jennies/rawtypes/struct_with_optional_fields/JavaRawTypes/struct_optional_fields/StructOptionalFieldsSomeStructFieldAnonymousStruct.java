package struct_optional_fields;

import java.util.Objects;

public class StructOptionalFieldsSomeStructFieldAnonymousStruct {
    public Object fieldAny;
    public StructOptionalFieldsSomeStructFieldAnonymousStruct() {
        this.fieldAny = new Object();
    }
    public StructOptionalFieldsSomeStructFieldAnonymousStruct(Object fieldAny) {
        this.fieldAny = fieldAny;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof StructOptionalFieldsSomeStructFieldAnonymousStruct)) return false;
        StructOptionalFieldsSomeStructFieldAnonymousStruct o = (StructOptionalFieldsSomeStructFieldAnonymousStruct) other;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldAny);
    }
}

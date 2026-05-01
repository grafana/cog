package struct_complex_fields;

import java.util.Objects;

public class StructComplexFieldsSomeStructFieldAnonymousStruct {
    public Object fieldAny;
    public StructComplexFieldsSomeStructFieldAnonymousStruct() {
        this.fieldAny = new Object();
    }
    public StructComplexFieldsSomeStructFieldAnonymousStruct(Object fieldAny) {
        this.fieldAny = fieldAny;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof StructComplexFieldsSomeStructFieldAnonymousStruct)) return false;
        StructComplexFieldsSomeStructFieldAnonymousStruct o = (StructComplexFieldsSomeStructFieldAnonymousStruct) other;
        if (!Objects.equals(this.fieldAny, o.fieldAny)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldAny);
    }
}

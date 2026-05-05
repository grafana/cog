package disjunctions_of_refs_without_discriminator;

import java.util.Objects;

public class TypeB {
    public Long fieldB;
    public TypeB() {
        this.fieldB = 0L;
    }
    public TypeB(Long fieldB) {
        this.fieldB = fieldB;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof TypeB)) return false;
        TypeB o = (TypeB) other;
        if (!Objects.equals(this.fieldB, o.fieldB)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldB);
    }
}

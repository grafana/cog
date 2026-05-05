package disjunctions_of_refs_without_discriminator;

import java.util.Objects;

public class TypeA {
    public String fieldA;
    public TypeA() {
        this.fieldA = "";
    }
    public TypeA(String fieldA) {
        this.fieldA = fieldA;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof TypeA)) return false;
        TypeA o = (TypeA) other;
        if (!Objects.equals(this.fieldA, o.fieldA)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldA);
    }
}

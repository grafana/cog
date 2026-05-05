package constant_references;

import java.util.Objects;

public class StructA {
    public Enum myEnum;
    public Enum other;
    public StructA() {
        this.myEnum = Enum.VALUE_A;
        this.other = Enum.VALUE_A;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof StructA)) return false;
        StructA o = (StructA) other;
        if (!Objects.equals(this.myEnum, o.myEnum)) return false;
        if (!Objects.equals(this.other, o.other)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.myEnum, this.other);
    }
}

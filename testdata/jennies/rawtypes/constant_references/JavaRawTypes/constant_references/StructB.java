package constant_references;

import java.util.Objects;

public class StructB {
    public Enum myEnum;
    public String myValue;
    public StructB() {
        this.myEnum = Enum.VALUE_B;
        this.myValue = "";
    }
    public StructB(String myValue) {
        this.myEnum = Enum.VALUE_B;
        this.myValue = myValue;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof StructB)) return false;
        StructB o = (StructB) other;
        if (!Objects.equals(this.myEnum, o.myEnum)) return false;
        if (!Objects.equals(this.myValue, o.myValue)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.myEnum, this.myValue);
    }
}

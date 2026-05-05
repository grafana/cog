package constant_references;

import java.util.Objects;

public class Struct {
    public String myValue;
    public Enum myEnum;
    public Struct() {
        this.myValue = "";
        this.myEnum = Enum.VALUE_A;
    }
    public Struct(String myValue,Enum myEnum) {
        this.myValue = myValue;
        this.myEnum = myEnum;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Struct)) return false;
        Struct o = (Struct) other;
        if (!Objects.equals(this.myValue, o.myValue)) return false;
        if (!Objects.equals(this.myEnum, o.myEnum)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.myValue, this.myEnum);
    }
}

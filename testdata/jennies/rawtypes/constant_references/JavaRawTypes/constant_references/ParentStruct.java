package constant_references;

import java.util.Objects;

public class ParentStruct {
    public Enum myEnum;
    public ParentStruct() {
        this.myEnum = Enum.VALUE_A;
    }
    public ParentStruct(Enum myEnum) {
        this.myEnum = myEnum;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof ParentStruct)) return false;
        ParentStruct o = (ParentStruct) other;
        if (!Objects.equals(this.myEnum, o.myEnum)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.myEnum);
    }
}

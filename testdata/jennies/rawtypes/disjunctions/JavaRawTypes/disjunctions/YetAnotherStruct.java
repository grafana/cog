package disjunctions;

import java.util.Objects;

public class YetAnotherStruct {
    public String type;
    public Integer bar;
    public YetAnotherStruct() {
        this.type = "";
        this.bar = 0;
    }
    public YetAnotherStruct(String type,Integer bar) {
        this.type = type;
        this.bar = bar;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof YetAnotherStruct)) return false;
        YetAnotherStruct o = (YetAnotherStruct) other;
        if (!Objects.equals(this.type, o.type)) return false;
        if (!Objects.equals(this.bar, o.bar)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.type, this.bar);
    }
}

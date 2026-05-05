package disjunctions_of_scalars_and_refs;

import java.util.Objects;

public class MyRefB {
    public Long bar;
    public MyRefB() {
        this.bar = 0L;
    }
    public MyRefB(Long bar) {
        this.bar = bar;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof MyRefB)) return false;
        MyRefB o = (MyRefB) other;
        if (!Objects.equals(this.bar, o.bar)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.bar);
    }
}

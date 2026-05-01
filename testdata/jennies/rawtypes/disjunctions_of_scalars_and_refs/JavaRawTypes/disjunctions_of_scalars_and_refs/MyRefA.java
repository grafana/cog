package disjunctions_of_scalars_and_refs;

import java.util.Objects;

public class MyRefA {
    public String foo;
    public MyRefA() {
        this.foo = "";
    }
    public MyRefA(String foo) {
        this.foo = foo;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof MyRefA)) return false;
        MyRefA o = (MyRefA) other;
        if (!Objects.equals(this.foo, o.foo)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.foo);
    }
}

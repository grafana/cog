package reference_of_reference;

import java.util.Objects;

public class OtherStruct {
    public String a;
    public OtherStruct() {
        this.a = "";
    }
    public OtherStruct(String a) {
        this.a = a;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof OtherStruct)) return false;
        OtherStruct o = (OtherStruct) other;
        if (!Objects.equals(this.a, o.a)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.a);
    }
}

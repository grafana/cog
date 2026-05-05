package nullable_fields;

import java.util.Objects;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Struct {
    public MyObject a;
    public MyObject b;
    public String c;
    public List<String> d;
    public Map<String, String> e;
    public NullableFieldsStructF f;
    public String g;
    public Struct() {
        this.e = new HashMap<>();
        this.g = Constants.ConstantRef;
    }
    public Struct(MyObject a,MyObject b,String c,List<String> d,Map<String, String> e,NullableFieldsStructF f) {
        this.a = a;
        this.b = b;
        this.c = c;
        this.d = d;
        this.e = e;
        this.f = f;
        this.g = Constants.ConstantRef;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Struct)) return false;
        Struct o = (Struct) other;
        if (!Objects.equals(this.a, o.a)) return false;
        if (!Objects.equals(this.b, o.b)) return false;
        if (!Objects.equals(this.c, o.c)) return false;
        if (!Objects.equals(this.d, o.d)) return false;
        if (!Objects.equals(this.e, o.e)) return false;
        if (!Objects.equals(this.f, o.f)) return false;
        if (!Objects.equals(this.g, o.g)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.a, this.b, this.c, this.d, this.e, this.f, this.g);
    }
}

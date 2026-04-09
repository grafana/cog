package nullable_fields;

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
}

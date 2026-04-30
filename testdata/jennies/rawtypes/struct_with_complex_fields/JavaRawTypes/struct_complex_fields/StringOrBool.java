package struct_complex_fields;

import java.util.Objects;

public class StringOrBool {
    protected String string;
    protected Boolean bool;
    protected StringOrBool() {}
    public static StringOrBool createString(String string) {
        StringOrBool stringOrBool = new StringOrBool();
        stringOrBool.string = string;
        return stringOrBool;
    }
    public static StringOrBool createBool(Boolean bool) {
        StringOrBool stringOrBool = new StringOrBool();
        stringOrBool.bool = bool;
        return stringOrBool;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof StringOrBool)) return false;
        StringOrBool o = (StringOrBool) other;
        if (!Objects.equals(this.string, o.string)) return false;
        if (!Objects.equals(this.bool, o.bool)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.string, this.bool);
    }
}

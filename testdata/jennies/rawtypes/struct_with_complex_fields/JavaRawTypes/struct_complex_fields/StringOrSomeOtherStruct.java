package struct_complex_fields;

import java.util.Objects;

public class StringOrSomeOtherStruct {
    protected String string;
    protected SomeOtherStruct someOtherStruct;
    protected StringOrSomeOtherStruct() {}
    public static StringOrSomeOtherStruct createString(String string) {
        StringOrSomeOtherStruct stringOrSomeOtherStruct = new StringOrSomeOtherStruct();
        stringOrSomeOtherStruct.string = string;
        return stringOrSomeOtherStruct;
    }
    public static StringOrSomeOtherStruct createSomeOtherStruct(SomeOtherStruct someOtherStruct) {
        StringOrSomeOtherStruct stringOrSomeOtherStruct = new StringOrSomeOtherStruct();
        stringOrSomeOtherStruct.someOtherStruct = someOtherStruct;
        return stringOrSomeOtherStruct;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof StringOrSomeOtherStruct)) return false;
        StringOrSomeOtherStruct o = (StringOrSomeOtherStruct) other;
        if (!Objects.equals(this.string, o.string)) return false;
        if (!Objects.equals(this.someOtherStruct, o.someOtherStruct)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.string, this.someOtherStruct);
    }
}

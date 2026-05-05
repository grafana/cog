package constant_reference_as_default;

import java.util.Objects;

public class MyStruct {
    public String aString;
    public String optString;
    public MyStruct() {
        this.aString = Constants.ConstantRefString;
        this.optString = Constants.ConstantRefString;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof MyStruct)) return false;
        MyStruct o = (MyStruct) other;
        if (!Objects.equals(this.aString, o.aString)) return false;
        if (!Objects.equals(this.optString, o.optString)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.aString, this.optString);
    }
}

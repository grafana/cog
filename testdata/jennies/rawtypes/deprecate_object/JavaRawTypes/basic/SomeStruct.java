package basic;

import java.util.Objects;

/** 
 * @deprecated This object is deprecated, use NewStruct instead.
 */
@Deprecated(forRemoval = true)
public class SomeStruct {
    public String fieldString;
    public SomeStruct() {
        this.fieldString = "";
    }
    public SomeStruct(String fieldString) {
        this.fieldString = fieldString;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.fieldString, o.fieldString)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.fieldString);
    }
}

package defaults;

import java.util.Objects;

public class DefaultsStructComplexFieldNested {
    public String nestedVal;
    public DefaultsStructComplexFieldNested() {
        this.nestedVal = "";
    }
    public DefaultsStructComplexFieldNested(String nestedVal) {
        this.nestedVal = nestedVal;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof DefaultsStructComplexFieldNested)) return false;
        DefaultsStructComplexFieldNested o = (DefaultsStructComplexFieldNested) other;
        if (!Objects.equals(this.nestedVal, o.nestedVal)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.nestedVal);
    }
}

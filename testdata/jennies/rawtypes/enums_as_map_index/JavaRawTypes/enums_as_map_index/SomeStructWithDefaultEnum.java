package enums_as_map_index;

import java.util.Objects;
import java.util.HashMap;
import java.util.Map;

public class SomeStructWithDefaultEnum {
    public Map<StringEnumWithDefault, String> data;
    public SomeStructWithDefaultEnum() {
        this.data = new HashMap<>();
    }
    public SomeStructWithDefaultEnum(Map<StringEnumWithDefault, String> data) {
        this.data = data;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStructWithDefaultEnum)) return false;
        SomeStructWithDefaultEnum o = (SomeStructWithDefaultEnum) other;
        if (!Objects.equals(this.data, o.data)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.data);
    }
}

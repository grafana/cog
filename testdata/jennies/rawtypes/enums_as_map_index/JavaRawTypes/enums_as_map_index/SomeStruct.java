package enums_as_map_index;

import java.util.Objects;
import java.util.HashMap;
import java.util.Map;

public class SomeStruct {
    public Map<StringEnum, String> data;
    public SomeStruct() {
        this.data = new HashMap<>();
    }
    public SomeStruct(Map<StringEnum, String> data) {
        this.data = data;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.data, o.data)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.data);
    }
}

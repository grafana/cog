package enums_as_map_index;

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
}

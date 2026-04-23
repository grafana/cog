package enums_as_map_index;

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
}

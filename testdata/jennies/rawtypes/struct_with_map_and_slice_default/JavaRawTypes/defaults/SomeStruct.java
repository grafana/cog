package defaults;

import java.util.Objects;
import java.util.Map;
import java.util.List;

public class SomeStruct {
    public Map<String, Object> options;
    public List<String> items;
    public Object extra;
    public SomeStruct() {
        this.options = Map.of();
        this.items = List.of();
        this.extra = Map.of();
    }
    public SomeStruct(Map<String, Object> options,List<String> items,Object extra) {
        this.options = options;
        this.items = items;
        this.extra = extra;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof SomeStruct)) return false;
        SomeStruct o = (SomeStruct) other;
        if (!Objects.equals(this.options, o.options)) return false;
        if (!Objects.equals(this.items, o.items)) return false;
        if (!Objects.equals(this.extra, o.extra)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.options, this.items, this.extra);
    }
}

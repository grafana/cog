package nullable_map_assignment;

import java.util.Map;

public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder config(Map<String, String> config) {
        this.internal.config = config;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}

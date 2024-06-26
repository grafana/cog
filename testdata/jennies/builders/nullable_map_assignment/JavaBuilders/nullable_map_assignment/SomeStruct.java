package nullable_map_assignment;

import java.util.Map;

public class SomeStruct {
    public Map<String, String> config;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Config(Map<String, String> config) {
    this.internal.config = config;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}

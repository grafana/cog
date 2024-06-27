package basic_struct_defaults;

import java.util.List;

public class SomeStruct {
    public Long id;
    public String uid;
    public List<String> tags;
    public Boolean liveNow;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        this.Id(42L);
        this.Uid("default-uid");
        this.Tags(List.of("generated", "cog"));
        this.LiveNow(true);
        }
    public Builder Id(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder Uid(String uid) {
    this.internal.uid = uid;
        return this;
    }
    
    public Builder Tags(List<String> tags) {
    this.internal.tags = tags;
        return this;
    }
    
    public Builder LiveNow(Boolean liveNow) {
    this.internal.liveNow = liveNow;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}

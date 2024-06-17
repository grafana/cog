package basic_struct_defaults;

import java.util.List;

public class SomeStruct {
    public Long id;
    public String uid;
    public List<String> tags;
    public Boolean liveNow;
        
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        this.setId(42L);
        this.setUid("default-uid");
        this.setTags(List.of("generated", "cog"));
        this.setLiveNow(true);
        }
    public Builder setId(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder setUid(String uid) {
    this.internal.uid = uid;
        return this;
    }
    
    public Builder setTags(List<String> tags) {
    this.internal.tags = tags;
        return this;
    }
    
    public Builder setLiveNow(Boolean liveNow) {
    this.internal.liveNow = liveNow;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}

package basic_struct;

import java.util.List;

// SomeStruct, to hold data.
public class SomeStruct {
    // id identifies something. Weird, right?
    public Long id;
    public String uid;
    public List<String> tags;
    // This thing could be live.
    // Or maybe not.
    public Boolean liveNow;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
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
    public SomeStruct build() {
            return this.internal;
        }
    }
}

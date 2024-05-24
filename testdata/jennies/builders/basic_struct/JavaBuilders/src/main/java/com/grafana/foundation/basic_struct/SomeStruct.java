package com.grafana.foundation.basic_struct;

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
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
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

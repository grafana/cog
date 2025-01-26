package basic_struct;

import java.util.List;

public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder id(Long id) {
        this.internal.id = id;
        return this;
    }
    
    public SomeStructBuilder uid(String uid) {
        this.internal.uid = uid;
        return this;
    }
    
    public SomeStructBuilder tags(List<String> tags) {
        this.internal.tags = tags;
        return this;
    }
    
    public SomeStructBuilder liveNow(Boolean liveNow) {
        this.internal.liveNow = liveNow;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}

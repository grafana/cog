package anonymous_struct;


public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder time(Object time) {
        this.internal.time = time;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}

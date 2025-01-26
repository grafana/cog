package sandbox;


public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder editable() {
        this.internal.editable = true;
        return this;
    }
    
    public SomeStructBuilder readonly() {
        this.internal.editable = false;
        return this;
    }
    
    public SomeStructBuilder autoRefresh() {
        this.internal.autoRefresh = true;
        return this;
    }
    
    public SomeStructBuilder noAutoRefresh() {
        this.internal.autoRefresh = false;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}

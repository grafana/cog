package struct_with_defaults;


public class NestedStructBuilder implements cog.Builder<NestedStruct> {
    protected final NestedStruct internal;
    
    public NestedStructBuilder() {
        this.internal = new NestedStruct();
    }
    public NestedStructBuilder stringVal(String stringVal) {
        this.internal.stringVal = stringVal;
        return this;
    }
    
    public NestedStructBuilder intVal(Long intVal) {
        this.internal.intVal = intVal;
        return this;
    }
    public NestedStruct build() {
        return this.internal;
    }
}

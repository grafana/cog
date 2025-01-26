package properties;


public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    private String someBuilderProperty;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    this.someBuilderProperty = "";
    }
    public SomeStructBuilder id(Long id) {
        this.internal.id = id;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}

package properties;


public class SomeStruct {
    public Long id;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        private String someBuilderProperty;
        
        public Builder() {
            this.internal = new SomeStruct();
        this.someBuilderProperty = "";
        }
    public Builder Id(Long id) {
    this.internal.id = id;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}

package properties;


public class SomeStruct {
    public Long id;
    
    public static class Builder {
        private SomeStruct internal;
        private String someBuilderProperty;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setId(Long id) {
    this.internal.id = id;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}

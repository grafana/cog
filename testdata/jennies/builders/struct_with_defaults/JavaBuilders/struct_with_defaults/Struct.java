package struct_with_defaults;


public class Struct {
    public NestedStruct allFields;
    public NestedStruct partialFields;
    public NestedStruct emptyFields;
    public Object complexField;
    public Object partialComplexField;
        
    public static class Builder {
        private Struct internal;
        
        public Builder() {
            this.internal = new Struct();
        NestedStruct.Builder nestedStructResource = new NestedStruct.Builder();
        nestedStructResource.setStringVal("hello");
        nestedStructResource.setIntVal(3L);
        this.setAllFields(nestedStructResource);
        NestedStruct.Builder nestedStructResource = new NestedStruct.Builder();
        nestedStructResource.setIntVal(4L);
        this.setPartialFields(nestedStructResource);
        this.setComplexField(new Object());
        this.setPartialComplexField(new Object());
        }
    public Builder setAllFields(cog.Builder<NestedStruct> allFields) {
    this.internal.allFields = allFields.build();
        return this;
    }
    
    public Builder setPartialFields(cog.Builder<NestedStruct> partialFields) {
    this.internal.partialFields = partialFields.build();
        return this;
    }
    
    public Builder setEmptyFields(cog.Builder<NestedStruct> emptyFields) {
    this.internal.emptyFields = emptyFields.build();
        return this;
    }
    
    public Builder setComplexField(Object complexField) {
    this.internal.complexField = complexField;
        return this;
    }
    
    public Builder setPartialComplexField(Object partialComplexField) {
    this.internal.partialComplexField = partialComplexField;
        return this;
    }
    public Struct build() {
            return this.internal;
        }
    }
}

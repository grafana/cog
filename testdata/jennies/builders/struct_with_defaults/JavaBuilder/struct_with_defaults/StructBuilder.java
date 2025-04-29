package struct_with_defaults;


public class StructBuilder implements cog.Builder<Struct> {
    protected final Struct internal;
    
    public StructBuilder() {
        this.internal = new Struct();
    }
    public StructBuilder allFields(cog.Builder<NestedStruct> allFields) {
    NestedStruct allFieldsResource = allFields.build();
        this.internal.allFields = allFieldsResource;
        return this;
    }
    
    public StructBuilder partialFields(cog.Builder<NestedStruct> partialFields) {
    NestedStruct partialFieldsResource = partialFields.build();
        this.internal.partialFields = partialFieldsResource;
        return this;
    }
    
    public StructBuilder emptyFields(cog.Builder<NestedStruct> emptyFields) {
    NestedStruct emptyFieldsResource = emptyFields.build();
        this.internal.emptyFields = emptyFieldsResource;
        return this;
    }
    
    public StructBuilder complexField(Object complexField) {
        this.internal.complexField = complexField;
        return this;
    }
    
    public StructBuilder partialComplexField(Object partialComplexField) {
        this.internal.partialComplexField = partialComplexField;
        return this;
    }
    public Struct build() {
        return this.internal;
    }
}

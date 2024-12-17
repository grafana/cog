package struct_with_defaults;

import struct_with_defaults.NestedStructBuilder;

public class StructBuilder implements cog.Builder<Struct> {
    protected final Struct internal;
    
    public StructBuilder() {
        this.internal = new Struct();
    NestedStructBuilder nestedStructResource = new NestedStructBuilder();
    nestedStructResource.stringVal("hello");
    nestedStructResource.intVal(3L);
    this.allFields(nestedStructResource);
    NestedStructBuilder nestedStructResource = new NestedStructBuilder();
    nestedStructResource.intVal(4L);
    this.partialFields(nestedStructResource);
    this.complexField(new Object());
    this.partialComplexField(new Object());
    }
    public StructBuilder allFields(cog.Builder<NestedStruct> allFields) {
    this.internal.allFields = allFields.build();
        return this;
    }
    
    public StructBuilder partialFields(cog.Builder<NestedStruct> partialFields) {
    this.internal.partialFields = partialFields.build();
        return this;
    }
    
    public StructBuilder emptyFields(cog.Builder<NestedStruct> emptyFields) {
    this.internal.emptyFields = emptyFields.build();
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

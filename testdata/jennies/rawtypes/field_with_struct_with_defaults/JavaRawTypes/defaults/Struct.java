package defaults;

import java.util.Objects;

public class Struct {
    public NestedStruct allFields;
    public NestedStruct partialFields;
    public NestedStruct emptyFields;
    public DefaultsStructComplexField complexField;
    public DefaultsStructPartialComplexField partialComplexField;
    public Struct() {
        this.allFields = new NestedStruct("hello", 3L);
        this.partialFields = new NestedStruct("", 3L);
        this.emptyFields = new defaults.NestedStruct();
        this.complexField = new DefaultsStructComplexField("myUID", new DefaultsStructComplexFieldNested("nested"), List.of("hello"));
        this.partialComplexField = new DefaultsStructPartialComplexField("", 0L);
    }
    public Struct(NestedStruct allFields,NestedStruct partialFields,NestedStruct emptyFields,DefaultsStructComplexField complexField,DefaultsStructPartialComplexField partialComplexField) {
        this.allFields = allFields;
        this.partialFields = partialFields;
        this.emptyFields = emptyFields;
        this.complexField = complexField;
        this.partialComplexField = partialComplexField;
    }

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Struct)) return false;
        Struct o = (Struct) other;
        if (!Objects.equals(this.allFields, o.allFields)) return false;
        if (!Objects.equals(this.partialFields, o.partialFields)) return false;
        if (!Objects.equals(this.emptyFields, o.emptyFields)) return false;
        if (!Objects.equals(this.complexField, o.complexField)) return false;
        if (!Objects.equals(this.partialComplexField, o.partialComplexField)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.allFields, this.partialFields, this.emptyFields, this.complexField, this.partialComplexField);
    }
}

package defaults;


public class Struct {
    public NestedStruct allFields;
    public NestedStruct partialFields;
    public NestedStruct emptyFields;
    public DefaultsStructComplexField complexField;
    public DefaultsStructPartialComplexField partialComplexField;
    public Struct() {}
    
    public Struct(NestedStruct allFields,NestedStruct partialFields,NestedStruct emptyFields,DefaultsStructComplexField complexField,DefaultsStructPartialComplexField partialComplexField) {
        this.allFields = allFields;
        this.partialFields = partialFields;
        this.emptyFields = emptyFields;
        this.complexField = complexField;
        this.partialComplexField = partialComplexField;
    }
}

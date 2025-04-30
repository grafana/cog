package defaults;


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
}

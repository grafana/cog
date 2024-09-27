package defaults;


public class Struct {
    public NestedStruct allFields;
    public NestedStruct partialFields;
    public NestedStruct emptyFields;
    public DefaultsStructComplexField complexField;
    public DefaultsStructPartialComplexField partialComplexField;

    public Struct() {
        NestedStruct nestedStructResource = new NestedStruct();
        nestedStructResource.stringVal = "hello";
        nestedStructResource.intVal = 3L;
        this.allFields = nestedStructResource;
        NestedStruct nestedStructResource = new NestedStruct();
        nestedStructResource.intVal = 3L;
        this.partialFields = nestedStructResource;
        DefaultsStructComplexField defaultsStructComplexFieldResource = new DefaultsStructComplexField();
        defaultsStructComplexFieldResource.uid = "myUID";
        this.complexField = defaultsStructComplexFieldResource;
        DefaultsStructPartialComplexField defaultsStructPartialComplexFieldResource = new DefaultsStructPartialComplexField();
        this.partialComplexField = defaultsStructPartialComplexFieldResource;
    }
}

namespace Grafana.Foundation.Defaults;


public class Struct
{
    public NestedStruct AllFields;
    public NestedStruct PartialFields;
    public NestedStruct EmptyFields;
    public DefaultsStructComplexField ComplexField;
    public DefaultsStructPartialComplexField PartialComplexField;

    public Struct()
    {
        this.AllFields = new NestedStruct();
        this.PartialFields = new NestedStruct();
        this.EmptyFields = new NestedStruct();
        this.ComplexField = new DefaultsStructComplexField();
        this.PartialComplexField = new DefaultsStructPartialComplexField();
    }

    public Struct(NestedStruct allFields, NestedStruct partialFields, NestedStruct emptyFields, DefaultsStructComplexField complexField, DefaultsStructPartialComplexField partialComplexField)
    {
        this.AllFields = allFields;
        this.PartialFields = partialFields;
        this.EmptyFields = emptyFields;
        this.ComplexField = complexField;
        this.PartialComplexField = partialComplexField;
    }
}
